package app

import (
	"app/internal/dto"
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func connectPg(databaseUrl string, attempt int8, ctx context.Context) (*pgxpool.Pool, error) {

	dbpool, err := pgxpool.New(ctx, databaseUrl)

	if err != nil {
		log.Println(err)
		time.Sleep(10 * time.Second)
		attempt--
		if attempt == 0 {
			return dbpool, err
		}
		connectPg(databaseUrl, attempt, ctx)
	}
	return dbpool, err
}

func (r *App) Run() {

	defer r.dbPool.Close()
	defer log.Println("Close DB connection pool")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go r.importListGorutine(wg)

	wg.Add(1)
	go r.importDetailsGorutine(wg)

	wg.Wait()
	log.Println("App shutdown")
}

// Горутина для асинхронного импорта списка сим-карт
func (r *App) importListGorutine(wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Shutdown importListGorutine")

	log.Println("Start importListGorutine")

	timer := time.NewTimer(0)

	for {

		if r.shutdown() {
			return
		}

		select {

		case <-timer.C:

			// Импорт списка сим-карт
			lostSimCards := r.importList()

			if lostSimCards != nil {
				r.createNotifications(lostSimCards)
				r.eraseLostSimCards(lostSimCards)
			}

			// Интевал между полными импортами списков сим-карт
			timer.Reset(10 * time.Minute)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Горутина для асинхронной работы сдетализацией сим-карт
func (r *App) importDetailsGorutine(wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Shutdown importDetailsGorutine")

	log.Println("Start importDetailsGorutine")

	timer := time.NewTimer(100 * time.Millisecond)

	for {

		if r.shutdown() {
			return
		}

		select {
		case <-timer.C:
			accountsList, err := r.spoRepo.GetOperatorsAccountsList()
			if err != nil {
				log.Println(err)
				return
			}
			r.getDetails(accountsList)
			timer.Reset(10 * time.Minute)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// Очистка списка заблокированных сим-карт
func (r *App) eraseLostSimCards(lostSimCards map[string]dto.LostSimCard) {
	for k := range lostSimCards {
		delete(lostSimCards, k)
	}
}

func (r *App) shutdown() bool {
	select {
	case <-r.ctx.Done():
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Printf("shudown handler %s\n", details.Name())
		}
		return true
	default:
		return false
	}
}
