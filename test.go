package main

import (
	"context"
	"errors"
	"fmt"
)

type User struct{ Name string }
type Game struct {
	players  []*User `json:"players"`
	gameName string  `json:"gameName"`
	number   int64   `json:"number"`
}

func main() {
	currentUser := &User{
		Name: "Jaehue",
	}

	// 컨텍스트 생성
	ctx := context.Background()

	// 컨텍스트에 값 추가 - context.WithValue 함수를 사용하여 새로운 컨텍스트를 생성함
	ctx = context.WithValue(ctx, "current_user", currentUser)

	// 함수 호출시 컨텍스트를 파라미터로 전달
	err := myFunc(ctx)
	if err != nil {
		println(err)
	}

	fmt.Println("============")
	var g Game
	g.gameName = "Witcher 3"
	g.printInformation()
	fmt.Println("addUser")
	g.addUser(currentUser)
	g.printInformation()
	fmt.Println("currentUser.Name")
	fmt.Println(currentUser.Name)
	g.changeUser("Jaehue", "HwangBo")
	fmt.Println("currentUser.Name")
	fmt.Println(currentUser.Name)
	g.printInformation()
}

func (g *Game) printInformation() {
	fmt.Println("Game Information")
	fmt.Println(g.gameName)
	for _, u := range g.players {
		fmt.Println(u.Name)
	}
	fmt.Println(g.number)
}

func (g *Game) addUser(user *User) {
	g.players = append(g.players, user)
	g.number += 1
}

func change(users []*User, from string, to string) error {
	for _, u := range users {
		if from == u.Name {
			u.Name = to
			return nil
		}
	}
	return nil
}

func (g *Game) changeUser(from string, to string) {
	change(g.players, from, to)
}

//context는 살짝 python의 딕셔너리 느낌
func myFunc(ctx context.Context) error {
	var currentUser *User

	// 컨텍스트에서 값을 가져옴
	if v := ctx.Value("current_user"); v != nil {
		u, ok := v.(*User)
		if !ok {
			return errors.New("Not authorized")
		}
		currentUser = u
	} else {
		return errors.New("Not authorized")
	}

	// currentUser를 사용하여 로직 처리
	fmt.Println(currentUser.Name)

	return nil
}
