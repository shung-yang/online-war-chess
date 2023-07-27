### Online War Chess
這款遊戲採多人連線戰棋玩法，玩家可以自訂戰棋的攻擊力、移動範圍等設定，不再受官方規則的限制，讓對戰更具彈性。目前，該遊戲仍處於開發階段（v1版本），我們計劃先開發一套官方的戰棋規則來進行演示。未來的開發方向詳情請參閱 [Wiki Page](https://github.com/shung-yang/online-war-chess/wiki/Future-Idea)

_online-war-chess 是遊戲的後端專案，前端部分則托管在 online_chess_frontend repository 中_

---

### Steps before you first run this project
* set .env according to .env.example
* run migration
`go run ./migration`

### StepS to  Start the app
1.  Create a .env file (please refer to .env.example for the file contents).
2.  `docker compose up`