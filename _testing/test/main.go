package main

import (
	"log"

	"github.com/casbin/casbin/v2/util"
)

func main() {
	// dao.SetupDatabase(&database.Database{
	// 	"mysql",
	// 	"root:catmaotu@tcp(127.0.0.1:3306)/goadmin?charset=utf8&parseTime=True&loc=Local&timeout=1000ms",
	// 	true,
	// })

	// err := trans.Exec(context.Background(), dao.DB, func(ctx context.Context) error {
	// 	err := models.CRole.BatchDelete(ctx, []int{2})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return models.CUser.BatchDelete(ctx, []int{3})
	// })
	// log.Println(err)

	log.Println(util.KeyMatch2("/abc/cde", "/:id"))
}
