# Go Fundamentals — คู่มือพื้นฐานภาษา Go (CS471 Day 2)

# 2. Building a Go API

เอกสารนี้เป็นเนื้อหาต่อเนื่องจาก Day 1 โดยตรง เนื้อหาจะเริ่มจากการรับ HTTP Request ด้วย Fiber แล้วค่อย ๆ แยกโค้ดออกเป็นชั้นตามหลัก Clean Architecture และจบด้วยการเปลี่ยนที่เก็บข้อมูลจาก Map ใน Memory ไปเป็น Postgres ที่รันด้วย Docker Compose

โปรเจกต์ที่ใช้ตลอดทั้งเอกสารคือ **ระบบสั่งอาหาร** ซึ่งเป็นตัวเดียวกับที่ใช้ใน Day 1

## 2.0 เตรียมเครื่องและโปรเจกต์ก่อนเริ่ม

### 2.0.1 สิ่งที่ต้องมีในเครื่อง

| สิ่งที่ต้องมี | ตรวจด้วยคำสั่ง | หมายเหตุ |
| --- | --- | --- |
| Go 1.22 ขึ้นไป | `go version` | ติดตั้งไว้แล้วตั้งแต่ Day 1 |
| Docker Desktop | `docker info` | ใช้ตั้งแต่หัวข้อ 2.5 เป็นต้นไป ต้องเปิดโปรแกรมค้างไว้ วิธีติดตั้งอยู่ใน `INSTALL-DOCKER.md` ผู้ที่มี Postgres ติดตั้งไว้ในเครื่องอยู่แล้ว ใช้ตัวนั้นแทนได้ ไม่ต้องติดตั้ง Docker |
| `curl` | `curl --version` | macOS และ Windows 10 ขึ้นไปมีมาให้อยู่แล้ว |

### 2.0.2 สร้างโฟลเดอร์ของวันนี้

วันนี้มี **2 โปรเจกต์แยกกัน** อยู่ในโฟลเดอร์เดียวกัน

<aside>
📄

**1. สร้างโฟลเดอร์หลัก** `mkdir week-02-go-api && cd week-02-go-api`

**2. สร้างโฟลเดอร์แรก** `mkdir fiber-playground && cd fiber-playground`

**3. สร้าง Module** `go mod init cs471/fiberplay`

**4. ติดตั้ง Fiber** `go get github.com/gofiber/fiber/v2`

**5. กลับขึ้นมาด้านบน** `cd ..`

</aside>

```bash
week-02-go-api/
├── fiber-playground/     # หัวข้อ 2.1 และ 2.2 ตัวอย่างสั้น ๆ แยกเป็นคนละโฟลเดอร์
└── food-api/             # หัวข้อ 2.4 และ 2.5 โปรเจกต์จริงที่ค่อย ๆ ขยายขึ้นทีละขั้น
```

<aside>
💡

`food-api` จะสร้างทีหลังที่ **ขั้นที่ 5** ในขั้นนี้ยังไม่ต้องสร้าง

เหตุผลที่แยกเป็นสอง Module เพราะ `fiber-playground` เป็นตัวอย่างที่ใช้แล้วทิ้งได้ ส่วน `food-api` คือโปรเจกต์จริงที่จะใช้ตลอดทั้งวัน

</aside>

### 2.0.3 วิธีอ่านเอกสารนี้

เนื้อหาตั้งแต่หัวข้อ 2.1 เป็นต้นไปแบ่งเป็น **37 ขั้น** เรียงตามลำดับที่จะสอนจริง ทุกขั้นมีกล่องสีระบุว่าต้องทำสิ่งใด และแบ่งเป็น **2 รูปแบบ** ที่แตกต่างกันชัดเจน

**รูปแบบที่ 1 — กล่องสีเทา หมายถึงการสร้างไฟล์ใหม่**

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p domain`

**2. สร้างไฟล์** `domain/menu.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**รูปแบบที่ 2 — กล่องสีส้ม หมายถึงการแก้ไฟล์เดิม**

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `usecase/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 7

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

ขั้นที่เป็นกล่องสีส้ม จะระบุข้อมูล 3 ส่วนเสมอ

1. **สิ่งที่แก้ไข** เป็นรายการโดยย่อว่ามีการแก้ไขสิ่งใดบ้าง
2. **โค้ดเฉพาะส่วนที่แก้ไข** สำหรับผู้ที่ต้องการแก้ไขด้วยตนเอง
3. **ไฟล์ฉบับเต็มหลังแก้ไข** อยู่ในกล่องที่กดเปิดได้ สำหรับผู้ที่เกิดความสับสนหรือแก้ไขผิดพลาด ให้ลบไฟล์เดิมแล้ววางฉบับนี้ทับลงไป

<aside>
🚨

หากทำตามแล้วเกิดความสับสน ให้กดเปิด **ไฟล์ฉบับเต็มหลังแก้ไข** ของขั้นนั้น แล้ววางทับทั้งไฟล์ เนื้อหาจะกลับมาตรงกับเอกสารทันที โดยไม่ต้องไล่ตรวจว่าผิดพลาดที่จุดใด

</aside>

### 2.0.4 คำสั่งที่ใช้บ่อย

| คำสั่ง | ใช้เมื่อใด |
| --- | --- |
| `go build ./...` | ตรวจสอบว่าทั้งโปรเจกต์ยังคอมไพล์ผ่าน โดยไม่ได้เรียกใช้โปรแกรม |
| `go run .` | รันเซิร์ฟเวอร์ จะทำงานค้างอยู่จนกว่าจะกด `Ctrl`  • `C` |
| `go test ./...` | รัน Test ทุกตัวในโปรเจกต์ |
| `go vet ./...` | ตรวจสอบจุดที่คอมไพล์ผ่านแต่มีแนวโน้มผิดพลาด เช่น Struct Tag พิมพ์ผิด |
| `gofmt -l .` | ระบุว่าไฟล์ใดยังจัดรูปแบบไม่ถูกต้อง หากไม่แสดงผลใด ๆ แสดงว่าเรียบร้อย |

<aside>
🚨

**ขั้นที่ต้องรันเซิร์ฟเวอร์ ต้องเปิด Terminal 2 หน้าต่าง**

หน้าต่างที่ 1 ปล่อย `go run .` ค้างไว้ ห้ามปิด

หน้าต่างที่ 2 ใช้เรียก `curl`

หากปรากฏข้อความ `address already in use` แสดงว่ามีโปรแกรมอื่นใช้ Port `3000` อยู่ ให้ปิดโปรแกรมนั้นก่อน หรือแก้เลข Port ใน `app.Listen` เป็น `3001` แล้วเปลี่ยนใน `curl` ตามด้วย

</aside>

## 2.1 Fiber

### 2.1.1 What a Framework Gives You

Fiber เป็น Web Framework ที่เขียนขึ้นบน `fasthttp` มีรูปแบบการเขียนใกล้เคียงกับ Express ใน Node.js สิ่งที่ Framework ช่วยจัดการให้มีดังนี้

1. **Routing** จับคู่ Method กับ Path ไปยังฟังก์ชันที่ต้องทำงาน พร้อมดึง Parameter ออกมาให้
2. **Middleware** งานอย่าง Logging, Authentication และ Recovery เขียนไว้ที่เดียวแล้วครอบทุก Request
3. **Binding** อ่าน JSON เข้า Struct และเขียน Struct กลับออกไปเป็น JSON
4. **Context** รวม Request, Response, Parameter และ State ไว้ในค่าเดียว
5. **Error Handling** มีจุดกลางที่แปลง Error ที่ถูกคืนกลับมาให้กลายเป็น Status Code

<aside>
💡

ทุกอย่างข้างต้นเขียนด้วย `net/http` ใน Standard Library ได้ทั้งหมด Fiber เป็นเพียงชั้นที่ทำให้เขียนสั้นลง ไม่ใช่สิ่งที่ขาดไม่ได้

</aside>

### 2.1.2 Hello, Fiber

#### ขั้นที่ 1 — เซิร์ฟเวอร์ Fiber ขนาดเล็กที่สุดที่ทำงานได้

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p 01-hello`

**2. สร้างไฟล์** `01-hello/main.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go run ./01-hello`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `fiber-playground`

</aside>

```go
// เซิร์ฟเวอร์ Fiber ขนาดเล็กที่สุดที่ทำงานได้
// รัน: go run ./01-hello   แล้วเปิด Terminal อีกหน้าต่างเพื่อเรียกใช้ curl

package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	// Handler ทุกตัวมีรูปแบบเดียวกัน คือ func(c *fiber.Ctx) error
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// บรรทัดนี้จะทำงานค้างอยู่ จนกว่าจะกด Ctrl+C
	app.Listen(":3000")
}
```

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

```bash
curl -i localhost:3000/health
```

```bash
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8

ok
```

- Handler ทุกตัวมีรูปแบบเดียวกันคือ `func(c *fiber.Ctx) error` การคืน Error ออกไปถือเป็นเรื่องปกติ ไม่ใช่ความผิดพลาด
- `*fiber.Ctx` เป็นค่าเดียวที่เก็บทั้ง Request และ Response ไว้ด้วยกัน
- `app.Listen` จะทำงานค้างอยู่จนกว่าโปรแกรมจะถูกปิด

<aside>
🚨

ห้ามเก็บ `*fiber.Ctx` ไว้ใช้ต่อหลังจาก Handler คืนค่าไปแล้ว เพราะ Fiber นำ Context กลับมาใช้ซ้ำเพื่อลดการจองหน่วยความจำ ค่าภายในจะถูกเขียนทับด้วย Request ของผู้ใช้รายอื่น

</aside>

### 2.1.3 Routing and Parameters

#### ขั้นที่ 2 — Path Parameter, Query Parameter และ Group

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p 02-routing`

**2. สร้างไฟล์** `02-routing/main.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go run ./02-routing`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `fiber-playground`

</aside>

```go
// Path Parameter, Query Parameter และ Group
// รัน: go run ./02-routing   แล้วเปิด Terminal อีกหน้าต่างเพื่อเรียกใช้ curl

package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	// Group ใช้รวม Prefix ไว้ที่เดียว หากต้องเปลี่ยนเป็น /api/v2 จะแก้เพียงบรรทัดเดียว
	api := app.Group("/api")

	api.Get("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")             // คืนค่าเป็น string เสมอ
		limit := c.QueryInt("limit", 20) // 20 คือค่าเริ่มต้น เมื่อผู้เรียกไม่ได้ส่งค่ามา

		return c.JSON(fiber.Map{"id": id, "limit": limit})
	})

	app.Listen(":3000")
}
```

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl localhost:3000/api/orders/o1
```

```bash
{"id":"o1","limit":20}
```

คำสั่งที่ 2

```bash
curl "localhost:3000/api/orders/o1?limit=5"
```

```bash
{"id":"o1","limit":5}
```

คำสั่งที่ 3

```bash
curl "localhost:3000/api/orders/o1?limit=abc"
```

```bash
{"id":"o1","limit":20}
```

- `c.Params` คืนค่าเป็น String เสมอ หากต้องการตัวเลข ต้องแปลงและตรวจสอบเอง
- `c.QueryInt` รับค่าเริ่มต้นเป็น Parameter ตัวที่สอง ใช้เมื่อผู้เรียกไม่ได้ส่งค่ามา
- `fiber.Map` เป็นชื่อเรียกแทน `map[string]any` ใช้เมื่อยังไม่จำเป็นต้องสร้าง Struct ขึ้นมาหนึ่งตัว

<aside>
💡

ใช้ `app.Group("/api")` เพื่อรวม Prefix ไว้ที่เดียว เวลาเปลี่ยนเป็น `/api/v2` จะแก้ไขเพียงบรรทัดเดียว

</aside>

### 2.1.4 Request and Response

#### ขั้นที่ 3 — อ่าน JSON เข้าสู่ Struct ด้วย BodyParser และเลือก Status Code

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p 03-body-parser`

**2. สร้างไฟล์** `03-body-parser/main.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go run ./03-body-parser`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `fiber-playground`

</aside>

```go
// อ่าน JSON เข้าสู่ Struct ด้วย BodyParser และเลือก Status Code ให้เหมาะสม
// รัน: go run ./03-body-parser   แล้วเปิด Terminal อีกหน้าต่างเพื่อเรียกใช้ curl

package main

import "github.com/gofiber/fiber/v2"

// ข้อความหลัง Backtick เรียกว่า Struct Tag
// ใช้ระบุว่าฟีลด์นี้ตรงกับ Key ชื่อใดใน JSON
type placeOrderRequest struct {
	MenuID string `json:"menu_id"`
	Qty    int    `json:"qty"`
}

func main() {
	app := fiber.New()

	app.Post("/api/orders", func(c *fiber.Ctx) error {
		var req placeOrderRequest

		// ต้องส่ง Pointer เข้าไป เพื่อให้เขียนค่าลงใน Struct ได้
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		if req.MenuID == "" || req.Qty < 1 {
			return c.Status(422).JSON(fiber.Map{"error": "menu_id and qty are required"})
		}

		return c.Status(201).JSON(req)
	})

	app.Listen(":3000")
}
```

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"m1","qty":2}'
```

```bash
HTTP/1.1 201 Created
Content-Type: application/json

{"menu_id":"m1","qty":2}
```

คำสั่งที่ 2

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d 'not json'
```

```bash
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error":"invalid body"}
```

คำสั่งที่ 3

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"","qty":0}'
```

```bash
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json

{"error":"menu_id and qty are required"}
```

- ข้อความหลัง Backtick เรียกว่า Struct Tag ใช้ระบุว่าฟีลด์นี้ตรงกับ Key ชื่อใดใน JSON
- `BodyParser` ต้องรับ Pointer จึงเขียนเป็น `&req` เพื่อให้เขียนค่าลงใน Struct ได้
- Struct ที่ใช้รับ Request ควรประกาศไว้ในชั้นนี้เท่านั้น ไม่ควรใช้ Struct ของ Domain โดยตรง

<aside>
🧠

**เจาะลึก** Struct Tag ไม่ใช่ Comment แต่เป็นข้อความที่ถูกเก็บไว้ในข้อมูลของ Type และอ่านออกมาได้ตอน Runtime ผ่าน Reflection

`encoding/json` จะอ่าน Tag เหล่านี้ตอนแปลงค่า ทำให้ชื่อฟีลด์ใน Go กับชื่อ Key ใน JSON ต่างกันได้ โดยไม่ต้องเขียนโค้ดแปลงเอง

Tag ที่เขียนผิดจะไม่มีการแจ้งเตือน เพราะคอมไพเลอร์ไม่ได้ตรวจสอบข้อความภายใน หากพิมพ์ `jsom` แทน `json` โปรแกรมจะยังคอมไพล์ผ่าน แต่ชื่อ Key จะกลายเป็นชื่อฟีลด์เดิม คำสั่ง `go vet` ช่วยตรวจจุดนี้ได้

</aside>

### 2.1.5 Status Codes

| Code | ความหมาย | ใช้เมื่อใด |
| --- | --- | --- |
| `200` | OK | ทำงานสำเร็จและมีข้อมูลส่งกลับไป |
| `201` | Created | สร้างข้อมูลใหม่สำเร็จ ควรส่งข้อมูลที่เพิ่งสร้างกลับไปด้วย |
| `204` | No Content | ทำงานสำเร็จแต่ไม่มีข้อมูลส่งกลับ เช่นการยกเลิกออเดอร์ |
| `400` | Bad Request | แปลง Body ไม่ได้ หรือฟีลด์ที่จำเป็นขาดหายไป |
| `404` | Not Found | ไม่พบเมนูหรือออเดอร์ที่ค้นหา |
| `409` | Conflict | คำขอถูกต้อง แต่สถานะปัจจุบันไม่อนุญาต เช่นยกเลิกตอนที่ร้านเริ่มทำอาหารแล้ว |
| `422` | Unprocessable Entity | แปลงค่าได้ แต่ผิดกฎทางธุรกิจ เช่นยอดไม่ถึงขั้นต่ำ |
| `500` | Internal Server Error | ระบบของเราเกิดข้อผิดพลาดเอง ห้ามส่งข้อความภายในกลับไปให้ผู้ใช้ |

<aside>
💡

ความต่างระหว่าง `400` กับ `422` คือ `400` หมายถึงอ่านคำขอไม่เข้าใจ ส่วน `422` หมายถึงอ่านเข้าใจแล้วแต่ทำตามไม่ได้เพราะผิดกฎ

</aside>

## 2.2 Methods and Errors in a Handler

### 2.2.1 Methods

เนื้อหาส่วนนี้เป็นการทบทวนจาก Day 1 เนื่องจากต้องนำไปใช้ทันทีในการเขียน Handler

```go
type Order struct {
	ID     string
	Lines  []OrderLine
	Status Status
}

// Value receiver: อ่านค่าอย่างเดียว
func (o Order) Total() float64 {
	sum := 0.0
	for _, l := range o.Lines {
		sum += l.Item.Price * float64(l.Qty)
	}
	return sum
}

// Pointer receiver: แก้ค่าต้นฉบับ
func (o *Order) Advance(next Status) {
	o.Status = next
}
```

### 2.2.2 Errors are Values

```go
// Sentinel Error ประกาศไว้ระดับ Package เพียงครั้งเดียว
var (
	ErrNotFound      = errors.New("not found")
	ErrBelowMinimum  = errors.New("order below 100 baht")
	ErrNotCancelable = errors.New("cooking already started")
)

// ห่อเพื่อเพิ่มบริบท โดยไม่ทำให้ต้นเหตุหายไป
if err != nil {
	return fmt.Errorf("place order %s: %w", id, err)
}

// ตรวจทะลุทุกชั้นที่ถูกห่อ
if errors.Is(err, ErrNotFound) {
	// 404
}
```

### 2.2.3 Mapping Errors to Status Codes

การตัดสินใจว่า Error แต่ละตัวควรแปลงเป็น Status Code ใด ควรเขียนไว้ที่เดียว แล้วให้ทุก Handler เรียกใช้

#### ขั้นที่ 4 — รวมการแปลง Error เป็น Status Code ไว้ที่เดียว

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p 04-error-mapping`

**2. สร้างไฟล์** `04-error-mapping/main.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go run ./04-error-mapping`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `fiber-playground`

</aside>

```go
// รวมการแปลง Error เป็น Status Code ไว้ที่เดียว
// รัน: go run ./04-error-mapping   แล้วเปิด Terminal อีกหน้าต่างเพื่อเรียกใช้ curl

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Sentinel Error ประกาศไว้ระดับ Package เพียงครั้งเดียว
var (
	ErrNotFound      = errors.New("not found")
	ErrBelowMinimum  = errors.New("order below 100 baht")
	ErrNotCancelable = errors.New("cooking already started")
)

func msg(s string) fiber.Map { return fiber.Map{"error": s} }

// จุดเดียวในโปรแกรมที่กำหนดว่า Error แต่ละตัวควรแปลงเป็น Status Code ใด
func toHTTP(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrNotFound):
		return c.Status(404).JSON(msg("not found"))

	case errors.Is(err, ErrBelowMinimum):
		return c.Status(422).JSON(msg(err.Error()))

	case errors.Is(err, ErrNotCancelable):
		return c.Status(409).JSON(msg(err.Error()))

	default:
		log.Println(err) // บันทึกไว้ในระบบ แต่ไม่ส่งรายละเอียดให้ผู้ใช้
		return c.Status(500).JSON(msg("internal error"))
	}
}

// จำลองชั้น Use Case ที่คืน Sentinel Error ออกมา โดยไม่รู้จัก HTTP
func fakeUseCase(kind string) error {
	switch kind {
	case "missing":
		return fmt.Errorf("place order %s: %w", kind, ErrNotFound)
	case "cheap":
		return ErrBelowMinimum
	case "cooking":
		return ErrNotCancelable
	case "boom":
		return errors.New("pq: relation \"orders\" does not exist")
	}
	return nil
}

func main() {
	app := fiber.New()

	app.Get("/try/:kind", func(c *fiber.Ctx) error {
		if err := fakeUseCase(c.Params("kind")); err != nil {
			return toHTTP(c, err)
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Listen(":3000")
}
```

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -i localhost:3000/try/ok
```

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

คำสั่งที่ 2

```bash
curl -i localhost:3000/try/missing
```

```bash
HTTP/1.1 404 Not Found
Content-Type: application/json

{"error":"not found"}
```

คำสั่งที่ 3

```bash
curl -i localhost:3000/try/cheap
```

```bash
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json

{"error":"order below 100 baht"}
```

คำสั่งที่ 4

```bash
curl -i localhost:3000/try/cooking
```

```bash
HTTP/1.1 409 Conflict
Content-Type: application/json

{"error":"cooking already started"}
```

คำสั่งที่ 5

```bash
curl -i localhost:3000/try/boom
```

```bash
HTTP/1.1 500 Internal Server Error
Content-Type: application/json

{"error":"internal error"}
```

- ชั้น Use Case เป็นผู้คืน Sentinel Error ออกมา โดยไม่รู้จัก HTTP
- มีเพียงชั้นนี้เท่านั้นที่ทราบว่า `ErrNotCancelable` ควรแปลงเป็น `409`
- Case สุดท้ายคือ `500` เสมอ เพราะ Error ที่ไม่รู้จักหมายความว่าเป็นข้อผิดพลาดของระบบเราเอง ให้บันทึกไว้ตรวจสอบ แต่ไม่ส่งให้ผู้ใช้

<aside>
🚨

สิ่งที่ไม่ควรทำในชั้น Handler

1. เขียน `c.Status(500)` กระจายอยู่ในทุกฟังก์ชัน
2. ส่ง `err.Error()` กลับไปให้ผู้ใช้โดยตรง เพราะอาจมีชื่อตารางหรือ Query ปะปนไปด้วย
3. เขียนกฎทางธุรกิจไว้ใน Handler
4. ละเลย Error ที่ได้จาก `BodyParser`
</aside>

## 2.3 Architecture

### 2.3.1 The Problem

โค้ดต่อไปนี้ทำงานได้ถูกต้องทุกประการ แต่เป็นตัวอย่างของสิ่งที่จะปรับปรุงตลอดทั้งบทนี้

```go
func PlaceOrder(c *fiber.Ctx) error {
	var req placeOrderRequest
	c.BodyParser(&req)

	// การคุยกับฐานข้อมูล อยู่ในชั้น Web
	row := db.QueryRow("SELECT price FROM menu WHERE id = $1", req.MenuID)
	var price float64
	if err := row.Scan(&price); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	// กฎทางธุรกิจ ฝังอยู่ใน Handler
	total := price * float64(req.Qty) + 30
	if total < 100 {
		return c.Status(422).JSON(fiber.Map{"error": "below minimum"})
	}

	db.Exec("INSERT INTO orders ...")
	return c.Status(201).JSON(fiber.Map{"total": total})
}
```

ปัญหาที่ตามมาในภายหลังมี 4 ข้อ

1. **ทดสอบไม่ได้** การทดสอบกฎ 100 บาท ต้องมีฐานข้อมูลที่ทำงานอยู่จริง และต้องส่ง HTTP Request เข้ามา
2. **ค้นหายาก** กฎยอดขั้นต่ำอยู่ใน Handler ซึ่งไม่มีผู้ใดคาดว่าจะพบกฎทางธุรกิจอยู่ในนั้น
3. **เปลี่ยนแปลงยาก** การเปลี่ยนจาก Postgres ไปใช้ระบบอื่น หมายความว่าต้องแก้ไข Handler ใหม่ทุกตัว
4. **เขียนซ้ำซ้อน** Endpoint ถัดไปที่ต้องใช้กฎเดียวกัน จะต้องเขียนกฎนั้นขึ้นใหม่อีกครั้ง

<aside>
💡

ให้จดจำโค้ดชุดนี้ไว้ เพราะที่ **ขั้นที่ 19** เราจะได้ Endpoint เดียวกันนี้ที่ทำงานเหมือนกันทุกอย่าง แต่แยกเป็น 4 ไฟล์ แล้วเปรียบเทียบว่าแตกต่างกันอย่างไร

</aside>

### 2.3.2 Four Ways People Organise Code

| รูปแบบ | แนวคิด |
| --- | --- |
| MVC | Model, View, Controller เป็นค่าเริ่มต้นของ Rails, Laravel และ Spring MVC เริ่มต้นง่ายที่สุด |
| Layered / N-tier | จัดวางเป็นชั้น Presentation อยู่บน Business และ Business อยู่บน Data โดยแต่ละชั้นเรียกใช้ชั้นที่อยู่ถัดลงไป |
| Hexagonal | เรียกอีกชื่อว่า Ports and Adapters แกนกลางเป็นผู้ประกาศ Port ส่วนภายนอกเป็นผู้เขียน Adapter มาเชื่อมต่อ |
| Clean | แนวคิดเดียวกับ Hexagonal แต่แบ่งเป็นวงชัดเจน และมีกฎข้อเดียวคือทิศทางการพึ่งพาชี้เข้าข้างใน |

### 2.3.3 MVC

MVC เป็นรูปแบบที่พบได้บ่อยที่สุดเมื่ออ่าน Tutorial ทั่วไป

1. **Model** ข้อมูล และบ่อยครั้งคือแถวในฐานข้อมูลโดยตรง
2. **View** Template หรือ JSON ที่ส่งกลับไป
3. **Controller** รับ Request ทำงาน แล้วตอบกลับ

<aside>
⚠️

ข้อจำกัดของ MVC คือไม่มีข้อกำหนดว่ากฎทางธุรกิจควรอยู่ที่ใด กฎจึงมารวมอยู่ใน Controller เมื่อเวลาผ่านไป Controller จะใหญ่ขึ้นเรื่อย ๆ ส่วน Model จะกลายเป็นเพียงแถวข้อมูลที่มี Logic ผูกติดอยู่

MVC จัดระเบียบ *การรับหนึ่ง Request* แต่ไม่ได้จัดระเบียบ *โดเมน* ช่องว่างนี้คือสิ่งที่ Clean Architecture พยายามปิด

</aside>

### 2.3.4 Clean Architecture

Clean Architecture แบ่งโค้ดออกเป็น 4 วง โดยวงในไม่รู้จักวงนอกเลย

| วง | โฟลเดอร์ของเรา | สิ่งที่อยู่ภายใน |
| --- | --- | --- |
| Entities | `domain/` | `Order`, `MenuItem`, `Status`, `CanCancel` กฎที่เป็นจริงแม้ไม่มีระบบนี้ |
| Use Cases | `usecase/` | `Place`, `Cancel`, `SetStatus` และ Interface ของ Repository ที่ตนต้องการ |
| Interface Adapters | `delivery/` `repository/` | Fiber Handler และส่วนที่ติดต่อกับ Postgres ทำหน้าที่แปลงข้อมูลทั้งสองทิศทาง |
| Frameworks and Drivers | `main.go` `pgx` `fiber` | ไลบรารีภายนอกและการประกอบทุกชั้นเข้าด้วยกัน เป็นที่เดียวที่อ้างถึงชื่อ Type จริง |

<aside>
💡

Clean, Hexagonal และ Onion คือแนวคิดเดียวกันที่เรียกชื่อต่างกัน ไม่จำเป็นต้องถกเถียงว่าชื่อใดถูกต้องกว่ากัน

</aside>

### 2.3.5 The Dependency Rule

กฎมีข้อเดียวคือ **การพึ่งพาในระดับ Source Code ชี้เข้าข้างในเท่านั้น** วงในไม่รู้จักสิ่งใดเลยเกี่ยวกับวงนอก

แผนภาพด้านล่างคือโครงของ `food-api` ที่จะสร้างในหัวข้อ 2.4 ลูกศรทุกเส้นชี้เข้าหาข้างใน ไม่มีเส้นใดชี้ออก

```mermaid
flowchart LR
    subgraph outer["Frameworks and Drivers"]
        MAIN["main.go"]
        FIBER["fiber"]
        PGX["pgx"]
    end
    subgraph adapters["Interface Adapters"]
        HTTP["delivery/http<br/>Handler, DTO, toHTTP"]
        MEM["repository/memory"]
        PG["repository/postgres"]
    end
    subgraph app["Use Cases"]
        UC["usecase<br/>OrderUseCase<br/>OrderRepository (interface)"]
    end
    subgraph ent["Entities"]
        DOM["domain<br/>Order, MenuItem, Status<br/>CanCancel, Total"]
    end

    MAIN --> HTTP
    MAIN --> MEM
    MAIN --> PG
    FIBER -.-> HTTP
    PGX -.-> PG
    HTTP --> UC
    MEM --> UC
    PG --> UC
    UC --> DOM
```

```go
// domain/order.go
import ()                       // ไม่ import อะไรเลย

// usecase/order.go
import "food-api/domain"        // ชี้เข้าข้างในเท่านั้น

// repository/postgres/order.go
import (
	"food-api/domain"           // เข้าข้างใน
	"github.com/jackc/pgx/v5"   // ชี้ออกข้างนอกได้ เพราะอยู่วงนอกอยู่แล้ว
)

// domain/order.go  <-- บรรทัดนี้คือการทำผิดกฎ
// import "github.com/jackc/pgx/v5"
```

- บล็อก `import` คือสถาปัตยกรรมที่มองเห็นได้ ตรวจสอบตอน Review ได้ในไม่กี่วินาที
- Go ช่วยบังคับให้ได้ส่วนหนึ่ง เพราะ Import Cycle เป็น Compile Error
- เมื่อวงในจำเป็นต้องใช้วงนอก ให้วงในประกาศ Interface แล้วให้วงนอกเขียนส่วนที่ทำงานจริงมาเชื่อมต่อ

### 2.3.6 Crossing Boundaries

ข้อมูลที่ส่งผ่านระหว่างวงต้องเป็นข้อมูลธรรมดา ไม่ใช่ค่าที่ผูกกับ Framework

```go
// delivery: รูปแบบที่ส่งออกไปภายนอก มี JSON Tag
type orderResponse struct {
	ID     string  `json:"id"`
	Total  float64 `json:"total"`
	Status string  `json:"status"`
}

// domain: ไม่มี Tag ไม่มี Framework ไม่มีฐานข้อมูล
type Order struct {
	ID     string
	Lines  []OrderLine
	Status Status
}

// Adapter ทำหน้าที่แปลงข้อมูลทางเดียว
func toResponse(o domain.Order) orderResponse { ... }
```

<aside>
🚨

ห้ามส่ง `*fiber.Ctx` หรือ `*sql.Rows` ข้ามวงเข้าไปข้างใน เพราะเท่ากับลาก Framework เข้าไปในกฎทางธุรกิจ

</aside>

### 2.3.7 The Names You Will Actually See

ตำราแต่ละเล่มเรียกชื่อทั้งสี่หน้าที่นี้ไม่เหมือนกัน ตารางนี้ใช้เทียบเคียงเมื่ออ่าน Tutorial จากแหล่งอื่น

| หน้าที่ | ในเอกสารนี้ | ที่อื่นเรียกว่า | เก็บสิ่งใด |
| --- | --- | --- | --- |
| Entity | `domain.Order` | Model, Entity | ข้อมูล และกฎที่เป็นจริงเสมอ |
| Use Case | `usecase.OrderUseCase` | Service, Interactor | หนึ่ง Method ต่อหนึ่งสิ่งที่ผู้ใช้กระทำได้ และเป็นผู้บังคับใช้กฎ |
| Repository | `memory.OrderRepo` | Repository, DAO, Store | อ่านและเขียน SQL ทั้งหมดอยู่ที่นี่ |
| Controller | `http.Handler` | Controller, Resource | แปลง Request เรียก Use Case หนึ่งตัว แล้วตอบกลับ |
| DTO | `orderResponse` | DTO, Request, Response | รูปแบบของข้อมูลที่ส่งผ่านข้ามขอบเขต |

<aside>
⚠️

หาก Tutorial เขียนว่า Service ให้เข้าใจว่าเป็น Use Case และหากเขียนว่า Model ให้ตรวจสอบก่อนว่าหมายถึงข้อมูลหรือหมายถึงแถวในฐานข้อมูล ความต่างตรงนี้คือแก่นของเรื่องทั้งหมด

</aside>

### 2.3.8 SRP and Dependency Inversion

- **Single Responsibility Principle** โมดูลหนึ่งควรมีเหตุผลให้ต้องแก้เพียงเหตุผลเดียว ไม่ใช่ทำงานเดียว แต่คือมีแหล่งของการเปลี่ยนแปลงแหล่งเดียว
- ในระบบของเรา รูปแบบของ API, กฎทางธุรกิจ และชนิดของฐานข้อมูล เปลี่ยนด้วยเหตุผลคนละอย่างกัน จึงไม่ควรอยู่ในไฟล์เดียวกัน
- **Dependency Inversion** ชั้นที่อยู่ข้างในเป็นผู้ประกาศว่าตนต้องการอะไร แล้วชั้นนอกปรับตัวมาให้ตรง ไม่ใช่ในทางกลับกัน

```go
package usecase

// ประกาศไว้ใน Package ที่เป็นผู้ใช้งาน
type OrderRepository interface {
	FindMenuItem(id string) (domain.MenuItem, error)
	Save(o domain.Order) error
}

type OrderUseCase struct {
	repo OrderRepository
}

func NewOrderUseCase(r OrderRepository) OrderUseCase {
	return OrderUseCase{repo: r}
}
```

ผลของการประกาศ Interface ไว้ฝั่งผู้ใช้งาน คือมี 3 ตัวที่ใช้แทนกันได้ทันที โดยที่ `OrderUseCase` ไม่รู้จักตัวใดเลย

```mermaid
classDiagram
    class OrderRepository {
        <<interface>>
        +ListMenu() []MenuItem, error
        +FindMenuItem(id) MenuItem, error
        +FindOrder(id) Order, error
        +Save(o Order) error
    }
    class OrderUseCase {
        -repo OrderRepository
        +ListMenu()
        +Place(menuID, qty)
        +Cancel(id)
        +SetStatus(id, next)
        +Total(id)
    }
    class MemoryOrderRepo {
        -menu map
        -orders map
    }
    class PostgresOrderRepo {
        -db pgxpool.Pool
    }
    class StubRepo {
        -order Order
    }

    OrderUseCase o-- OrderRepository : ถือไว้เป็นฟีลด์
    OrderRepository <|.. MemoryOrderRepo : implements
    OrderRepository <|.. PostgresOrderRepo : implements
    OrderRepository <|.. StubRepo : implements (ใช้ใน Test)
```

<aside>
💡

หลักที่ควรจำคือ **รับ Interface เข้ามาเป็น Parameter แต่คืน Struct ตัวจริงออกไป** เพราะผู้เรียกควรได้รับสิ่งที่ใช้งานได้เต็มที่ ส่วนเราควรระบุความต้องการให้แคบที่สุด

</aside>

## 2.4 Building the Order API

### 2.4.1 The Brief

ก่อนเขียนโค้ดบรรทัดแรก ต้องทราบก่อนว่ากำลังสร้างสิ่งใด

**สิ่งที่ระบบต้องทำได้**

1. ลูกค้าดูเมนูได้
2. ลูกค้าสั่งอาหารได้
3. ลูกค้าเปิดดูออเดอร์จาก `id` ได้
4. ร้านเปลี่ยนสถานะของออเดอร์ได้
5. ลูกค้ายกเลิกได้ ก่อนที่ร้านจะเริ่มทำอาหาร

ค่าคงที่ของระบบคือ ยอดสั่งซื้อขั้นต่ำ 100 บาท และค่าจัดส่ง 30 บาท

**สิ่งที่ยังไม่ทำในวันนี้** ได้แก่ ระบบบัญชีผู้ใช้และสิทธิ์การเข้าถึง, การชำระเงินและการคืนเงิน, การรองรับหลายร้าน, การจ่ายงานให้ไรเดอร์และการติดตามแบบเรียลไทม์ และการแจ้งเตือน

<aside>
💡

การระบุว่า *สิ่งใดจะไม่ทำ* สำคัญไม่น้อยกว่าการระบุว่าสิ่งใดต้องทำ เพราะเป็นสิ่งที่ป้องกันไม่ให้ขอบเขตงานขยายเกินความจำเป็นระหว่างทาง

</aside>

### 2.4.2 One Request, End to End

เส้นทางของ `POST /api/orders` ตั้งแต่รับเข้ามาจนตอบกลับ

```mermaid
sequenceDiagram
    autonumber
    actor C as Client
    participant H as delivery/http<br/>Handler
    participant U as usecase<br/>OrderUseCase
    participant R as repository<br/>OrderRepo
    participant D as domain<br/>Order

    C->>H: POST /api/orders {menu_id, qty}
    H->>H: BodyParser เข้า placeOrderRequest
    H->>U: Place(menuID, qty)
    U->>R: FindMenuItem(id)
    alt ไม่พบเมนู
        R-->>U: domain.ErrNotFound
        U-->>H: place order: not found
        H-->>C: 404 Not Found
    else พบเมนู
        R-->>U: domain.MenuItem
        U->>D: order.Total()
        D-->>U: 160
        alt ต่ำกว่า 100 บาท
            U-->>H: domain.ErrBelowMinimum
            H-->>C: 422 Unprocessable Entity
        else ผ่านกฎ
            U->>R: Save(order)
            R-->>U: nil
            U-->>H: domain.Order
            H->>H: toResponse(order)
            H-->>C: 201 Created + JSON
        end
    end
```

| ลำดับ | ไฟล์ | เกิดสิ่งใดขึ้น |
| --- | --- | --- |
| 1 | `delivery/http/order.go` | `BodyParser` อ่าน JSON เข้า `placeOrderRequest` ไม่มีกฎทางธุรกิจอยู่ในชั้นนี้ |
| 2 | `usecase/order.go` | `Place()` ขอเมนูจาก Repository แล้วบังคับกฎ 100 บาท |
| 3 | `repository/memory/order.go` | ค้นหาเมนูและบันทึกออเดอร์ ไม่มีกฎทางธุรกิจอยู่ในชั้นนี้เช่นกัน |
| 4 | `domain/order.go` | `Total()` และ `CanCancel()` ตอบคำถามที่ Use Case ถาม |
| 5 | `delivery/http/order.go` | `toResponse` สร้าง JSON แล้วส่ง `201 Created` กลับออกไป |

### 2.4.3 Project Structure

ในภาษา Go หนึ่งโฟลเดอร์คือหนึ่ง Package โครงสร้างโฟลเดอร์จึงเป็นสถาปัตยกรรมโดยตรง นี่คือโครงสร้างของ `food-api` เมื่อทำครบทุกขั้นแล้ว

```bash
food-api/
├── main.go                        # ประกอบทุกชั้นเข้าด้วยกัน
├── go.mod
├── .env.example                   # ขั้นที่ 28  ไฟล์ตัวอย่าง ลง Git ได้
├── .env                           # ขั้นที่ 28  ค่าจริงของเครื่องเรา ห้ามลง Git
├── .gitignore                     # ขั้นที่ 28
├── docker-compose.yml             # ขั้นที่ 29  ใช้เฉพาะผู้ที่เปิด Postgres ด้วย Docker
├── cmd/
│   └── migrate/main.go            # ขั้นที่ 32  ตัวรัน Migration
├── migrations/
│   ├── 001_init.sql               # ขั้นที่ 30
│   ├── 002_add_note_to_orders.sql # ขั้นที่ 33
│   └── 003_index_orders_status.sql# ขั้นที่ 34
├── domain/
│   ├── menu.go                    # ขั้นที่ 6
│   ├── order.go                   # ขั้นที่ 12
│   ├── errors.go                  # ขั้นที่ 13
│   └── order_test.go              # ขั้นที่ 14
├── usecase/
│   ├── order.go                   # ขั้นที่ 7  แล้วโตขึ้นที่ขั้น 15, 20, 24
│   ├── fakes_test.go              # ขั้นที่ 27
│   └── order_test.go              # ขั้นที่ 27
├── repository/
│   ├── memory/order.go            # ขั้นที่ 8  แล้วโตขึ้นที่ขั้น 15
│   └── postgres/order.go          # ขั้นที่ 35
└── delivery/
    └── http/
        ├── dto.go                 # ขั้นที่ 9  แล้วโตขึ้นที่ขั้น 17, 21
        ├── error.go               # ขั้นที่ 16
        └── order.go               # ขั้นที่ 10 แล้วโตขึ้นที่ขั้น 18, 22, 25
```

<aside>
💡

เริ่มจากสี่โฟลเดอร์นี้ก็เพียงพอ ค่อยเพิ่มระดับความลึกเมื่อชั้นใดชั้นหนึ่งขยายใหญ่ขึ้นจริง ไม่จำเป็นต้องสร้างโครงสร้างที่ซับซ้อนไว้ล่วงหน้า

</aside>

### 2.4.4 เริ่มโปรเจกต์

#### ขั้นที่ 5 — สร้าง Module ของโปรเจกต์ และติดตั้ง Fiber

<aside>
📄

**1. สร้างโฟลเดอร์แล้วเข้าไปข้างใน** `mkdir food-api && cd food-api`

**2. สร้าง Module** `go mod init food-api`

**3. ติดตั้ง Fiber** `go get github.com/gofiber/fiber/v2`

</aside>

```go
module food-api

go 1.22
```

- บรรทัด `go` แสดงรุ่นของ Go ที่ติดตั้งอยู่ในเครื่อง จึงอาจขึ้นเป็นเลขอื่นที่สูงกว่าตัวอย่างนี้ ไม่ต้องแก้ให้ตรงกัน

<aside>
🚨

ชื่อ Module ต้องเป็น `food-api` เป๊ะ ๆ เพราะทุก `import` ในเอกสารนี้ขึ้นต้นด้วยชื่อนี้ เช่น `food-api/domain` หากตั้งชื่ออื่น จะต้องแก้ไข `import` ทุกไฟล์เอง

</aside>

### 2.4.5 ชั้นที่ 1 — Domain

วงในสุด ไม่รู้จักชั้นใดเลย เริ่มจากสิ่งที่เล็กที่สุด คือรายการอาหารหนึ่งรายการ

#### ขั้นที่ 6 — ชั้นที่ 1 Domain — MenuItem ที่ไม่รู้จักทั้ง HTTP และฐานข้อมูล

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p domain`

**2. สร้างไฟล์** `domain/menu.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package domain

// MenuItem คือรายการอาหารหนึ่งรายการ
//
// สังเกตว่าไม่มี JSON Tag ไม่มี db Tag และไม่ import สิ่งใดเลย
// ชั้นนี้ไม่รู้จักทั้ง HTTP และฐานข้อมูล
type MenuItem struct {
	ID    string
	Name  string
	Price float64
}
```

- ไม่มี Fiber ไม่มี SQL ไม่มี JSON Tag ไฟล์นี้ทดสอบได้โดยไม่ต้องเตรียมสิ่งใดเลย
- ยังไม่ต้องสร้าง `Order` ในขั้นนี้ เป้าหมายของ 6 ขั้นแรก คือทำให้ `GET /api/menu` วิ่งผ่านครบทั้ง 4 ชั้นก่อน แล้วจึงค่อยเพิ่มส่วนที่ซับซ้อนขึ้น

### 2.4.6 ชั้นที่ 2 — Use Case

ชั้นนี้เป็นผู้ประกาศว่า *ตนต้องการสิ่งใดจากชั้นที่เก็บข้อมูล* โดยไม่จำเป็นต้องทราบว่าผู้ใดจะเป็นผู้ทำให้

#### ขั้นที่ 7 — ชั้นที่ 2 Use Case — ประกาศ Interface ของสิ่งที่ตัวเองต้องการ

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p usecase`

**2. สร้างไฟล์** `usecase/order.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package usecase

import (
	"food-api/domain"
)

// Interface ประกาศไว้ใน Package ที่เป็นผู้ใช้งาน ไม่ใช่ที่ผู้ให้บริการ
// นี่คือ Dependency Inversion ทิศทางการพึ่งพาจึงชี้เข้าข้างใน
type OrderRepository interface {
	ListMenu() ([]domain.MenuItem, error)
}

type OrderUseCase struct {
	repo OrderRepository
}

// รับ Interface เข้ามาเป็นพารามิเตอร์ แต่คืน Struct ตัวจริงออกไป
func NewOrderUseCase(r OrderRepository) OrderUseCase {
	return OrderUseCase{repo: r}
}

func (u OrderUseCase) ListMenu() ([]domain.MenuItem, error) {
	return u.repo.ListMenu()
}
```

- `OrderRepository` ประกาศอยู่ใน Package `usecase` ไม่ใช่ใน Package ของฐานข้อมูล นี่คือหัวใจของ Dependency Inversion
- ในขั้นนี้ Interface มี Method เดียวก็เพียงพอ แล้วค่อยขยายเมื่อจำเป็นต้องใช้จริง

### 2.4.7 ชั้นที่ 3 — Repository

ส่วนที่ทำงานจริงตาม Interface ด้านบน โดยรอบนี้จัดเก็บไว้ใน Map ธรรมดา

#### ขั้นที่ 8 — ชั้นที่ 3 Repository — เก็บข้อมูลใน Memory

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p repository/memory`

**2. สร้างไฟล์** `repository/memory/order.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package memory

import (
	"sort"
	"sync"

	"food-api/domain"
)

// ชั้นนี้ทำหน้าที่จัดเก็บและอ่านข้อมูลเท่านั้น ไม่มีกฎทางธุรกิจ
type OrderRepo struct {
	mu   *sync.RWMutex
	menu map[string]domain.MenuItem
}

func NewOrderRepo(menu map[string]domain.MenuItem) OrderRepo {
	return OrderRepo{
		mu:   &sync.RWMutex{},
		menu: menu,
	}
}

func (r OrderRepo) ListMenu() ([]domain.MenuItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]domain.MenuItem, 0, len(r.menu))
	for _, m := range r.menu {
		items = append(items, m)
	}

	// ลำดับของ range บน map ถูกสุ่มใหม่ทุกครั้ง (Day 1)
	// จึงต้องเรียงลำดับเอง มิฉะนั้น API จะคืนลำดับไม่เหมือนเดิมในแต่ละครั้ง
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}
```

- ชั้นนี้จัดเก็บและอ่านข้อมูลเท่านั้น ไม่มีกฎทางธุรกิจ
- `sort.Slice` จำเป็น เพราะลำดับของ `range` บน map ถูกสุ่มใหม่ทุกครั้ง หากไม่เรียงลำดับ API จะคืนลำดับไม่เหมือนเดิมในแต่ละครั้ง
- `OrderRepo` ไม่ได้เขียนว่า `implements OrderRepository` ที่ไหนเลย เพียงมี Method ครบถ้วน Go ก็ถือว่าใช้แทนกันได้ทันที

<aside>
🧠

**คำถาม เหตุใดจึงต้องใช้ `defer` ทั้งที่ภาษา Go มี Garbage Collector อยู่แล้ว**

Garbage Collector ดูแล **หน่วยความจำ** เพียงอย่างเดียว หน้าที่ของมันคือคืนหน่วยความจำของค่าที่ไม่มีผู้ใดอ้างถึงอีกต่อไป และจะคืนเมื่อใดก็ได้ที่เห็นสมควร โดยไม่รับประกันว่าจะเกิดขึ้นตอนใด

สิ่งที่ `defer` ดูแลไม่ใช่หน่วยความจำ แต่คือ **สิ่งที่ยืมมาแล้วต้องคืน** ซึ่ง Garbage Collector มองไม่เห็นเลย ได้แก่ Lock ของ Mutex, ไฟล์ที่เปิดค้างไว้, การเชื่อมต่อฐานข้อมูล, Transaction ที่เปิดไว้ และ Body ของ HTTP Response

ในโค้ดด้านบน `r.mu.RLock()` ไม่ได้จองหน่วยความจำแม้แต่ไบต์เดียว แต่เปลี่ยน **สถานะ** ของโปรแกรม หากลืมเรียก `RUnlock` สิ่งที่เกิดขึ้นไม่ใช่หน่วยความจำรั่ว แต่คือ Goroutine ตัวถัดไปที่ขอ Lock จะรอค้างอยู่ตลอดไป ตัว `mu` ยังมีผู้อ้างถึงและยังถูกใช้งานอยู่ Garbage Collector จึงไม่มีเหตุผลใดที่จะเข้ามายุ่ง และช่วยอะไรไม่ได้เลย

เหตุผลที่ใช้ `defer` แทนการเขียน `RUnlock()` ไว้บรรทัดสุดท้าย คือ `defer` ทำงานบน **ทุกเส้นทางที่ออกจากฟังก์ชัน** ทั้ง `return` ทุกจุด และตอนเกิด Panic เมื่อฟังก์ชันมีหลาย `return` การเขียนเองจะหลงลืมได้ง่าย

ภาษา Go ไม่มี Destructor แบบ C++ และไม่มี `try-with-resources` แบบ Java จึงใช้ `defer` ทำหน้าที่นี้แทน

สรุปคือ Garbage Collector จัดการหน่วยความจำ ส่วน `defer` จัดการสิ่งที่ต้องคืนตามเวลาที่แน่นอน ทั้งสองทำงานคนละหน้าที่กัน และใช้แทนกันไม่ได้

</aside>

### 2.4.8 ชั้นที่ 4 — Delivery

ชั้นนอกสุดของโค้ดเรา แบ่งเป็นสองไฟล์ คือรูปแบบของข้อมูล และตัว Handler

#### ขั้นที่ 9 — ชั้นที่ 4 Delivery — DTO ที่แยกจาก Domain

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p delivery/http`

**2. สร้างไฟล์** `delivery/http/dto.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package http

import (
	"food-api/domain"

	"github.com/gofiber/fiber/v2"
)

// รูปแบบของข้อมูลที่ส่งออกไป แยกจาก Type ใน domain
// จึงเปลี่ยนชื่อ Key ใน JSON ได้โดยไม่ต้องแก้ไข Domain
type menuResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Adapter ทำหน้าที่แปลงข้อมูลทางเดียว จาก domain ออกสู่ภายนอก
func toMenuResponse(items []domain.MenuItem) []menuResponse {
	out := make([]menuResponse, 0, len(items))
	for _, m := range items {
		out = append(out, menuResponse{ID: m.ID, Name: m.Name, Price: m.Price})
	}
	return out
}

func msg(s string) fiber.Map { return fiber.Map{"error": s} }
```

- `menuResponse` มี JSON Tag ส่วน `domain.MenuItem` ไม่มี ทั้งสองตัวมีข้อมูลเหมือนกันแต่คนละหน้าที่
- เปลี่ยนชื่อ Key ที่ส่งออกไปเป็น `menu_id` ได้ โดยไม่ต้องแก้ไข `domain` เลย

#### ขั้นที่ 10 — Handler ตัวแรก ที่เรียก Use Case ไม่ใช่เรียกฐานข้อมูล

<aside>
📄

**1. สร้างไฟล์** `delivery/http/order.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package http

import (
	"food-api/usecase"

	"github.com/gofiber/fiber/v2"
)

// Handler ถือ Use Case ไว้ ไม่ได้ถือการเชื่อมต่อฐานข้อมูล
type Handler struct {
	uc usecase.OrderUseCase
}

func NewHandler(uc usecase.OrderUseCase) Handler {
	return Handler{uc: uc}
}

// ชั้นนี้ทำเพียงสามอย่าง คือแปลง Request เรียก Use Case หนึ่งตัว แล้วเขียน Response
func (h Handler) ListMenu(c *fiber.Ctx) error {
	items, err := h.uc.ListMenu()
	if err != nil {
		return c.Status(500).JSON(msg("internal error"))
	}
	return c.JSON(toMenuResponse(items))
}
```

- Handler ถือ `usecase.OrderUseCase` ไว้ ไม่ได้ถือ Repository และไม่ได้ถือการเชื่อมต่อฐานข้อมูล
- ยังไม่มี `toHTTP` ที่ขั้นนี้ จึงเขียน `c.Status(500)` ไปก่อน แล้วจะแก้ที่ **ขั้นที่ 18**

### 2.4.9 ประกอบทุกชั้นเข้าด้วยกันใน main.go

ถึงตรงนี้มีครบทั้ง 4 ชั้นแล้ว เหลือเพียงการประกอบเข้าด้วยกัน

#### ขั้นที่ 11 — ประกอบทุกชั้นเข้าด้วยกันใน main.go แล้วเรียกใช้ API ครั้งแรก

<aside>
📄

**1. สร้างไฟล์** `main.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package main

import (
	"log"

	"food-api/delivery/http"
	"food-api/domain"
	"food-api/repository/memory"
	"food-api/usecase"

	"github.com/gofiber/fiber/v2"
)

func seedMenu() map[string]domain.MenuItem {
	return map[string]domain.MenuItem{
		"m1": {ID: "m1", Name: "Pad Thai", Price: 80},
		"m2": {ID: "m2", Name: "Tom Yum Goong", Price: 120},
		"m3": {ID: "m3", Name: "Green Curry", Price: 95},
		"m4": {ID: "m4", Name: "Thai Iced Tea", Price: 45},
	}
}

// เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึง memory, usecase และ http พร้อมกัน
// ประกอบจากชั้นในออกมาสู่ชั้นนอก
func main() {
	// 1. ที่เก็บข้อมูล
	repo := memory.NewOrderRepo(seedMenu())

	// 2. กฎ
	uc := usecase.NewOrderUseCase(repo)

	// 3. ชั้นขนส่ง
	h := http.NewHandler(uc)

	app := fiber.New()

	api := app.Group("/api")
	api.Get("/menu", h.ListMenu)

	log.Fatal(app.Listen(":3000"))
}
```

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

```bash
curl -i localhost:3000/api/menu
```

```bash
HTTP/1.1 200 OK
Content-Type: application/json

[{"id":"m1","name":"Pad Thai","price":80},{"id":"m2","name":"Tom Yum Goong","price":120},{"id":"m3","name":"Green Curry","price":95},{"id":"m4","name":"Thai Iced Tea","price":45}]
```

<aside>
💡

`main.go` เป็นไฟล์เดียวที่อ้างถึง `memory`, `usecase` และ `http` พร้อมกัน โดยประกอบจากชั้นในออกมาสู่ชั้นนอก คือสร้างที่เก็บข้อมูลก่อน ส่งเข้าไปให้ชั้นกฎ แล้วจึงส่งชั้นกฎเข้าไปให้ชั้นขนส่ง

หากต้องการเปลี่ยนไปใช้ Postgres ก็แก้ไขเพียงบรรทัดเดียวในไฟล์นี้ ซึ่งจะทำจริงที่ **ขั้นที่ 36**

</aside>

<aside>
🚨

หากรันแล้วปรากฏข้อความ `address already in use` แสดงว่ามีโปรแกรมอื่นใช้ Port `3000` อยู่แล้ว ให้ปิดโปรแกรมนั้น หรือแก้ `app.Listen(":3000")` เป็น `":3001"` แล้วเปลี่ยนเลขใน `curl` ตามด้วย

</aside>

### 2.4.10 เพิ่ม Order เข้าไปใน Domain

โครง 4 ชั้นใช้ได้แล้ว ลำดับถัดไปคือการเพิ่มเนื้อหาจริง โดยเริ่มจากวงในสุดเสมอ

#### ขั้นที่ 12 — เพิ่ม Order, Status และกฎการยกเลิกเข้าไปใน Domain

<aside>
📄

**1. สร้างไฟล์** `domain/order.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package domain

// ค่าคงที่ทางธุรกิจ อยู่ในชั้นในสุด เพราะเป็นจริงเสมอไม่ว่าจะเก็บข้อมูลด้วยวิธีใด
const (
	MinimumOrder = 100.0 // ยอดสั่งซื้อขั้นต่ำ ไม่รวมค่าจัดส่ง
	DeliveryFee  = 30.0  // ค่าจัดส่ง
)

// Status ประกาศเป็น Type ของตัวเอง ไม่ใช่ string ธรรมดา
// คอมไพเลอร์จึงช่วยตรวจจับการพิมพ์ผิดได้
type Status string

const (
	Pending    Status = "pending"
	Confirmed  Status = "confirmed"
	Cooking    Status = "cooking"
	Delivering Status = "delivering"
	Completed  Status = "completed"
	Cancelled  Status = "cancelled"
)

// ลำดับของสถานะ เดินหน้าได้ทางเดียว
var flow = []Status{Pending, Confirmed, Cooking, Delivering, Completed}

type OrderLine struct {
	Item MenuItem
	Qty  int
}

func (l OrderLine) Total() float64 {
	return l.Item.Price * float64(l.Qty)
}

type Order struct {
	ID     string
	Lines  []OrderLine
	Status Status
}

// Value receiver เพราะอ่านอย่างเดียว
func (o Order) Total() float64 {
	sum := 0.0
	for _, l := range o.Lines {
		sum += l.Total()
	}
	return sum
}

// กฎอยู่ติดกับ Type ที่กฎนั้นควบคุม ไม่ใช่อยู่ใน Handler
func (o Order) CanCancel() bool {
	return o.Status == Pending || o.Status == Confirmed
}

// เดินหน้าได้ทีละขั้นเท่านั้น ไม่อนุญาตให้ข้ามขั้นหรือย้อนกลับ
func (o Order) CanAdvanceTo(next Status) bool {
	for i, s := range flow {
		if s == o.Status {
			return i+1 < len(flow) && flow[i+1] == next
		}
	}
	return false // สถานะ cancelled หรือ completed เดินต่อไม่ได้อีก
}
```

- `type Status string` ดีกว่าการใช้ `string` ธรรมดา เพราะคอมไพเลอร์ช่วยตรวจจับการพิมพ์ผิดได้
- `CanCancel` และ `CanAdvanceTo` ควรอยู่ข้าง `Order` ไม่ใช่อยู่ใน Handler
- `MinimumOrder` และ `DeliveryFee` อยู่ในชั้นนี้ เพราะเป็นกฎของธุรกิจ ไม่ใช่เรื่องของ HTTP หรือฐานข้อมูล

#### ขั้นที่ 13 — Sentinel Error ของ Domain

<aside>
📄

**1. สร้างไฟล์** `domain/errors.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package domain

import "errors"

// Sentinel Error ประกาศไว้ระดับ Package เพียงครั้งเดียว
// ชั้นในเป็นผู้คืนค่าออกไป ส่วนชั้นนอกเป็นผู้แปลงเป็น Status Code
var (
	ErrNotFound      = errors.New("not found")
	ErrBelowMinimum  = errors.New("order below 100 baht")
	ErrNotCancelable = errors.New("cooking already started")
	ErrInvalidStatus = errors.New("invalid status transition")
	ErrInvalidQty    = errors.New("qty must be at least 1")
)
```

- Error ทั้ง 5 ตัวประกาศไว้ที่เดียว แล้วทุกชั้นอ้างถึงตัวเดียวกันนี้
- ชั้นในเป็นผู้คืนค่าออกไป ส่วนชั้นนอกเป็นผู้ตัดสินว่าแต่ละตัวควรแปลงเป็น Status Code ใด ซึ่งจะทำที่ **ขั้นที่ 16**

### 2.4.11 The Order Lifecycle

สถานะของออเดอร์เดินหน้าทางเดียว และมีทางออกเพียงทางเดียวคือการยกเลิก

```mermaid
stateDiagram-v2
    [*] --> pending: POST /api/orders
    pending --> confirmed: PATCH status
    confirmed --> cooking: PATCH status
    cooking --> delivering: PATCH status
    delivering --> completed: PATCH status
    completed --> [*]

    pending --> cancelled: POST cancel (204)
    confirmed --> cancelled: POST cancel (204)
    cancelled --> [*]

    note right of confirmed
        ยกเลิกได้ถึงตรงนี้
        CanCancel() = true
    end note

    note right of cooking
        เริ่มทำอาหารแล้ว
        ยกเลิกไม่ได้ คืน 409
    end note
```

| ช่วง | การยกเลิก |
| --- | --- |
| `pending` และ `confirmed` | ยกเลิกได้ เพราะยังไม่มีการลงมือทำอาหาร ออเดอร์เพียงหยุดลง |
| `cooking` เป็นต้นไป | ยกเลิกไม่ได้ ต้องคืน `409 Conflict` เพราะอาหารถูกทำไปแล้ว |

<aside>
💡

กฎ *ยกเลิกไม่ได้เมื่อเริ่มทำอาหารแล้ว* คือเหตุผลที่ชัดเจนที่สุดว่าเหตุใดจึงต้องมีชั้น Use Case เพราะไม่มี Constraint ใดในฐานข้อมูลที่เขียนกฎข้อนี้ได้

</aside>

กฎเหล่านี้พิสูจน์ได้ด้วย Test ที่ไม่ต้องใช้ฐานข้อมูลและไม่ต้องเรียก HTTP เลย

#### ขั้นที่ 14 — ทดสอบกฎทางธุรกิจ โดยไม่ต้องมีฐานข้อมูลและไม่ต้องเรียก HTTP

<aside>
📄

**1. สร้างไฟล์** `domain/order_test.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go test ./domain/ -v`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package domain_test

import (
	"testing"

	"food-api/domain"
)

// ชั้น domain ทดสอบได้โดยไม่ต้องมีฐานข้อมูลและไม่ต้องมีเซิร์ฟเวอร์
// ไม่ต้องเตรียมสิ่งใดเลยนอกจากคำสั่ง go test
func TestCanCancel(t *testing.T) {
	cases := []struct {
		status domain.Status
		want   bool
	}{
		{domain.Pending, true},
		{domain.Confirmed, true},
		{domain.Cooking, false},
		{domain.Delivering, false},
		{domain.Completed, false},
		{domain.Cancelled, false},
	}

	for _, c := range cases {
		got := domain.Order{Status: c.status}.CanCancel()
		if got != c.want {
			t.Errorf("สถานะ %s: ได้ %v ต้องการ %v", c.status, got, c.want)
		}
	}
}

func TestTotal(t *testing.T) {
	o := domain.Order{Lines: []domain.OrderLine{
		{Item: domain.MenuItem{Name: "Pad Thai", Price: 80}, Qty: 2},
		{Item: domain.MenuItem{Name: "Tom Yum", Price: 120}, Qty: 1},
	}}

	if got := o.Total(); got != 280 {
		t.Errorf("ต้องการ 280 แต่ได้ %v", got)
	}
}

func TestCanAdvanceTo(t *testing.T) {
	if !(domain.Order{Status: domain.Pending}).CanAdvanceTo(domain.Confirmed) {
		t.Error("pending → confirmed ควรเดินหน้าได้")
	}
	if (domain.Order{Status: domain.Pending}).CanAdvanceTo(domain.Cooking) {
		t.Error("pending → cooking เป็นการข้ามขั้น ไม่ควรเดินหน้าได้")
	}
	if (domain.Order{Status: domain.Cooking}).CanAdvanceTo(domain.Confirmed) {
		t.Error("cooking → confirmed เป็นการย้อนกลับ ไม่ควรเดินหน้าได้")
	}
	if (domain.Order{Status: domain.Cancelled}).CanAdvanceTo(domain.Cooking) {
		t.Error("ออเดอร์ที่ยกเลิกแล้ว ต้องเดินหน้าต่อไม่ได้")
	}
}
```

**ผลลัพธ์**

```bash
go test ./domain/ -v
```

```bash
=== RUN   TestCanCancel
--- PASS: TestCanCancel (0.00s)
=== RUN   TestTotal
--- PASS: TestTotal (0.00s)
=== RUN   TestCanAdvanceTo
--- PASS: TestCanAdvanceTo (0.00s)
PASS
ok  	food-api/domain	(cached)
```

<aside>
💡

Test ชุดนี้ทำงานเสร็จในเสี้ยววินาที และจะยังรันได้เหมือนเดิมแม้จะเปลี่ยนไปใช้ Postgres ที่หัวข้อ 2.5 เพราะไม่ได้ผูกกับฐานข้อมูลตั้งแต่แรก

</aside>

### 2.4.12 เพิ่มการสั่งอาหาร

ตอนนี้ `usecase` ต้องการข้อมูลจาก Repository มากขึ้น จึงต้องขยาย Interface และทำให้ Repository ครบตามนั้น

<aside>
🚨

ขั้นนี้มี **2 ไฟล์ที่ต้องแก้** และต้องแก้ให้ครบทั้งคู่ก่อนรัน `go build`

หากแก้ไขเพียง `usecase/order.go` ไฟล์เดียว แล้วสั่ง Build จะปรากฏ Error ดังนี้

`memory.OrderRepo does not implement usecase.OrderRepository (missing method FindMenuItem)`

แสดงว่า Interface ขยายไปแล้ว แต่ Repository ยังไม่ได้แก้ไขตาม ให้แก้ไฟล์ที่ 2 ให้เสร็จก่อน

</aside>

#### ขั้นที่ 15 — เพิ่มการสั่งอาหาร — ขยาย Interface และทำให้ Repository ครบตามนั้น

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `usecase/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 7

**2. เปิดไฟล์** `repository/memory/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 8

**3. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**4. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `usecase/order.go`**

**สิ่งที่แก้ไข**

- ขยาย `OrderRepository` จาก Method เดียว เป็น 4 Method
- เพิ่มตัวนับหมายเลขออเดอร์ `seq` และ Method `nextID()`
- เพิ่ม Method `Place` ที่บังคับกฎยอดขั้นต่ำ และ Method `Get`
- เพิ่ม `fmt` กับ `sync/atomic` เข้าไปในบล็อก `import`

แก้ไขบล็อก `import` ที่ส่วนหัวของไฟล์ ให้เป็นดังนี้

```go
import (
	"fmt"
	"sync/atomic"

	"food-api/domain"
)
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `usecase/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package usecase
    
    import (
    	"fmt"
    	"sync/atomic"
    
    	"food-api/domain"
    )
    
    // Interface ประกาศไว้ใน Package ที่เป็นผู้ใช้งาน ไม่ใช่ที่ผู้ให้บริการ
    // นี่คือ Dependency Inversion ทิศทางการพึ่งพาจึงชี้เข้าข้างใน
    type OrderRepository interface {
    	ListMenu() ([]domain.MenuItem, error)
    	FindMenuItem(id string) (domain.MenuItem, error)
    	FindOrder(id string) (domain.Order, error)
    	Save(o domain.Order) error
    }
    
    type OrderUseCase struct {
    	repo OrderRepository
    	seq  *int64
    }
    
    // รับ Interface เข้ามาเป็นพารามิเตอร์ แต่คืน Struct ตัวจริงออกไป
    func NewOrderUseCase(r OrderRepository) OrderUseCase {
    	var n int64
    	return OrderUseCase{repo: r, seq: &n}
    }
    
    func (u OrderUseCase) nextID() string {
    	return fmt.Sprintf("o%d", atomic.AddInt64(u.seq, 1))
    }
    
    func (u OrderUseCase) ListMenu() ([]domain.MenuItem, error) {
    	return u.repo.ListMenu()
    }
    
    // หนึ่ง Method ต่อหนึ่งสิ่งที่ผู้ใช้กระทำได้
    func (u OrderUseCase) Place(menuID string, qty int) (domain.Order, error) {
    	if qty < 1 {
    		return domain.Order{}, domain.ErrInvalidQty
    	}
    
    	item, err := u.repo.FindMenuItem(menuID)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    
    	order := domain.Order{
    		Lines:  []domain.OrderLine{{Item: item, Qty: qty}},
    		Status: domain.Pending,
    	}
    
    	// ตรวจสอบกฎให้ครบก่อน แล้วจึงกำหนด ID และบันทึก
    	// หากกำหนด ID ตั้งแต่ต้น ออเดอร์ที่ถูกปฏิเสธจะใช้หมายเลขสิ้นเปลืองโดยเปล่าประโยชน์
    	if order.Total() < domain.MinimumOrder {
    		return domain.Order{}, domain.ErrBelowMinimum
    	}
    
    	order.ID = u.nextID()
    
    	if err := u.repo.Save(order); err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    	return order, nil
    }
    
    func (u OrderUseCase) Get(id string) (domain.Order, error) {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("get order %s: %w", id, err)
    	}
    	return order, nil
    }
    ```
    

**ไฟล์ที่ 2 — `repository/memory/order.go`**

**สิ่งที่แก้ไข**

- เพิ่มฟีลด์ `orders` เข้าไปใน Struct และสร้าง map ให้ใน `NewOrderRepo`
- เพิ่ม Method `FindMenuItem`, `FindOrder` และ `Save` ให้ครบตาม Interface ที่ขยายไป

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `repository/memory/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package memory
    
    import (
    	"sort"
    	"sync"
    
    	"food-api/domain"
    )
    
    // ชั้นนี้ทำหน้าที่จัดเก็บและอ่านข้อมูลเท่านั้น
    // กฎยอดขั้นต่ำ 100 บาท ไม่ใช่หน้าที่ของชั้นนี้
    type OrderRepo struct {
    	mu     *sync.RWMutex
    	menu   map[string]domain.MenuItem
    	orders map[string]domain.Order
    }
    
    func NewOrderRepo(menu map[string]domain.MenuItem) OrderRepo {
    	return OrderRepo{
    		mu:     &sync.RWMutex{},
    		menu:   menu,
    		orders: map[string]domain.Order{},
    	}
    }
    
    func (r OrderRepo) ListMenu() ([]domain.MenuItem, error) {
    	r.mu.RLock()
    	defer r.mu.RUnlock()
    
    	items := make([]domain.MenuItem, 0, len(r.menu))
    	for _, m := range r.menu {
    		items = append(items, m)
    	}
    
    	// ลำดับของ range บน map ถูกสุ่มใหม่ทุกครั้ง (Day 1)
    	// จึงต้องเรียงลำดับเอง มิฉะนั้น API จะคืนลำดับไม่เหมือนเดิมในแต่ละครั้ง
    	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
    	return items, nil
    }
    
    func (r OrderRepo) FindMenuItem(id string) (domain.MenuItem, error) {
    	r.mu.RLock()
    	defer r.mu.RUnlock()
    
    	item, ok := r.menu[id] // comma-ok ใช้แยกกรณี "ไม่มี key" ออกจาก "มี key แต่ค่าเป็น zero value"
    	if !ok {
    		// คืน Error ของ Domain ไม่ใช่ Error ของชั้นที่เก็บข้อมูล
    		return domain.MenuItem{}, domain.ErrNotFound
    	}
    	return item, nil
    }
    
    func (r OrderRepo) FindOrder(id string) (domain.Order, error) {
    	r.mu.RLock()
    	defer r.mu.RUnlock()
    
    	order, ok := r.orders[id]
    	if !ok {
    		return domain.Order{}, domain.ErrNotFound
    	}
    	return order, nil
    }
    
    // ใช้ Value receiver แต่ยังแก้ไขข้อมูลได้ เพราะ map เป็น Reference Type (Day 1)
    func (r OrderRepo) Save(o domain.Order) error {
    	r.mu.Lock()
    	defer r.mu.Unlock()
    
    	r.orders[o.ID] = o
    	return nil
    }
    ```
    
- `Place` ตรวจสอบกฎให้ครบก่อน แล้วจึงกำหนด ID และบันทึก ไม่ใช่บันทึกก่อนแล้วจึงตรวจสอบ
- ลำดับนี้สำคัญ หากกำหนด ID ตั้งแต่ต้น ออเดอร์ที่ถูกปฏิเสธจะใช้หมายเลขสิ้นเปลืองโดยเปล่าประโยชน์ ทำให้ออเดอร์แรกที่สำเร็จจริงไม่ได้หมายเลข `o1`
- `Place` คืน Sentinel Error ออกไป แล้วให้ชั้น Delivery เป็นผู้แปลงเป็น `422`

### 2.4.13 แปลง Error เป็น Status Code

ตอนนี้ `usecase` คืน Error ได้หลายแบบแล้ว ถึงเวลาสร้างจุดกลางที่แปลงเป็น Status Code

#### ขั้นที่ 16 — จุดเดียวที่แปลง Error เป็น Status Code

<aside>
📄

**1. สร้างไฟล์** `delivery/http/error.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package http

import (
	"errors"
	"log"

	"food-api/domain"

	"github.com/gofiber/fiber/v2"
)

// จุดเดียวในโปรแกรมที่กำหนดว่า Error แต่ละตัวควรแปลงเป็น Status Code ใด
// ชั้น Use Case คืน Sentinel Error ออกมา โดยไม่รู้จัก HTTP
func toHTTP(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.Status(404).JSON(msg("not found"))

	case errors.Is(err, domain.ErrBelowMinimum):
		return c.Status(422).JSON(msg(err.Error()))

	case errors.Is(err, domain.ErrNotCancelable):
		return c.Status(409).JSON(msg(err.Error()))

	case errors.Is(err, domain.ErrInvalidStatus):
		return c.Status(409).JSON(msg(err.Error()))

	case errors.Is(err, domain.ErrInvalidQty):
		return c.Status(422).JSON(msg(err.Error()))

	default:
		// Error ที่ไม่รู้จัก หมายความว่าเป็นข้อผิดพลาดของระบบเราเอง
		// บันทึกไว้ตรวจสอบภายหลัง แต่ไม่ส่งให้ผู้ใช้ เพราะอาจมีชื่อตารางปะปนอยู่
		log.Println(err)
		return c.Status(500).JSON(msg("internal error"))
	}
}
```

- ไฟล์นี้คือ **จุดเดียวในโปรแกรม** ที่กำหนดว่า Error แต่ละตัวควรแปลงเป็น Status Code ใด
- `default` เป็น `500` เสมอ เพราะ Error ที่ไม่รู้จักหมายความว่าเป็นข้อผิดพลาดของระบบเราเอง

#### ขั้นที่ 17 — เพิ่ม DTO ของออเดอร์

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `delivery/http/dto.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 9

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม `placeOrderRequest` สำหรับรับ JSON เข้ามา
- เพิ่ม `orderResponse` และ `lineResp` สำหรับส่งออกไป
- เพิ่มฟังก์ชัน `toResponse` ที่แปลง `domain.Order` เป็น DTO

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `delivery/http/dto.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/domain"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // รูปแบบของข้อมูลที่รับเข้ามา มี JSON Tag เพราะเป็นเรื่องของชั้นนอก
    type placeOrderRequest struct {
    	MenuID string `json:"menu_id"`
    	Qty    int    `json:"qty"`
    }
    
    // รูปแบบของข้อมูลที่ส่งออกไป แยกจาก Type ใน domain
    // จึงเปลี่ยนชื่อ Key ใน JSON ได้โดยไม่ต้องแก้ไข Domain
    type orderResponse struct {
    	ID     string     `json:"id"`
    	Status string     `json:"status"`
    	Total  float64    `json:"total"`
    	Lines  []lineResp `json:"lines"`
    }
    
    type lineResp struct {
    	MenuID string  `json:"menu_id"`
    	Name   string  `json:"name"`
    	Price  float64 `json:"price"`
    	Qty    int     `json:"qty"`
    }
    
    type menuResponse struct {
    	ID    string  `json:"id"`
    	Name  string  `json:"name"`
    	Price float64 `json:"price"`
    }
    
    // Adapter ทำหน้าที่แปลงข้อมูลทางเดียว จาก domain ออกสู่ภายนอก
    func toResponse(o domain.Order) orderResponse {
    	lines := make([]lineResp, 0, len(o.Lines))
    	for _, l := range o.Lines {
    		lines = append(lines, lineResp{
    			MenuID: l.Item.ID,
    			Name:   l.Item.Name,
    			Price:  l.Item.Price,
    			Qty:    l.Qty,
    		})
    	}
    
    	return orderResponse{
    		ID:     o.ID,
    		Status: string(o.Status),
    		Total:  o.Total(),
    		Lines:  lines,
    	}
    }
    
    func toMenuResponse(items []domain.MenuItem) []menuResponse {
    	out := make([]menuResponse, 0, len(items))
    	for _, m := range items {
    		out = append(out, menuResponse{ID: m.ID, Name: m.Name, Price: m.Price})
    	}
    	return out
    }
    
    func msg(s string) fiber.Map { return fiber.Map{"error": s} }
    ```
    

#### ขั้นที่ 18 — Handler ของการสั่งอาหารและการเปิดดูออเดอร์

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `delivery/http/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 10

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เปลี่ยน `ListMenu` ให้เรียก `toHTTP` แทนการเขียน `c.Status(500)` เอง
- เพิ่ม Handler `PlaceOrder` และ `GetOrder`

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `delivery/http/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // Handler ถือ Use Case ไว้ ไม่ได้ถือการเชื่อมต่อฐานข้อมูล
    type Handler struct {
    	uc usecase.OrderUseCase
    }
    
    func NewHandler(uc usecase.OrderUseCase) Handler {
    	return Handler{uc: uc}
    }
    
    // ชั้นนี้ทำเพียงสามอย่าง คือแปลง Request เรียก Use Case หนึ่งตัว แล้วเขียน Response
    // ไม่ปรากฏเลข 100 อยู่ในไฟล์นี้แม้แต่แห่งเดียว
    func (h Handler) ListMenu(c *fiber.Ctx) error {
    	items, err := h.uc.ListMenu()
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toMenuResponse(items))
    }
    
    func (h Handler) PlaceOrder(c *fiber.Ctx) error {
    	var req placeOrderRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	order, err := h.uc.Place(req.MenuID, req.Qty)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	return c.Status(201).JSON(toResponse(order))
    }
    
    func (h Handler) GetOrder(c *fiber.Ctx) error {
    	order, err := h.uc.Get(c.Params("id"))
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toResponse(order))
    }
    ```
    

### 2.4.14 สั่งอาหารได้จริง

เหลือขั้นตอนสุดท้าย คือเปิด Route ให้เข้าถึงได้

#### ขั้นที่ 19 — เพิ่ม Route แล้วสั่งอาหารได้จริง

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 11

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม Middleware ที่ Log ทุก Request ไว้ที่เดียว
- เพิ่ม Route `POST /api/orders` และ `GET /api/orders/:id`

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"log"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    func seedMenu() map[string]domain.MenuItem {
    	return map[string]domain.MenuItem{
    		"m1": {ID: "m1", Name: "Pad Thai", Price: 80},
    		"m2": {ID: "m2", Name: "Tom Yum Goong", Price: 120},
    		"m3": {ID: "m3", Name: "Green Curry", Price: 95},
    		"m4": {ID: "m4", Name: "Thai Iced Tea", Price: 45},
    	}
    }
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึง memory, usecase และ http พร้อมกัน
    // ประกอบจากชั้นในออกมาสู่ชั้นนอก
    func main() {
    	// 1. ที่เก็บข้อมูล
    	repo := memory.NewOrderRepo(seedMenu())
    
    	// 2. กฎ
    	uc := usecase.NewOrderUseCase(repo)
    
    	// 3. ชั้นขนส่ง
    	h := http.NewHandler(uc)
    
    	app := fiber.New()
    
    	// Middleware ครอบทุก Request โดยเขียนไว้ที่เดียว (Decorator Pattern)
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	api := app.Group("/api")
    	api.Get("/menu", h.ListMenu)
    	api.Post("/orders", h.PlaceOrder)
    	api.Get("/orders/:id", h.GetOrder)
    
    	log.Fatal(app.Listen(":3000"))
    }
    ```
    

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"m1","qty":2}'
```

```bash
HTTP/1.1 201 Created
Content-Type: application/json

{"id":"o1","status":"pending","total":160,"lines":[{"menu_id":"m1","name":"Pad Thai","price":80,"qty":2}]}
```

คำสั่งที่ 2

```bash
curl -i localhost:3000/api/orders/o1
```

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{"id":"o1","status":"pending","total":160,"lines":[{"menu_id":"m1","name":"Pad Thai","price":80,"qty":2}]}
```

คำสั่งที่ 3

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"m4","qty":1}'
```

```bash
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json

{"error":"order below 100 baht"}
```

คำสั่งที่ 4

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"m9","qty":1}'
```

```bash
HTTP/1.1 404 Not Found
Content-Type: application/json

{"error":"not found"}
```

คำสั่งที่ 5

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d 'not json'
```

```bash
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error":"invalid body"}
```

<aside>
💡

เทียบกับโค้ดในหัวข้อ 2.3.1 ที่ทุกอย่างรวมอยู่ในฟังก์ชันเดียว ตอนนี้ได้ผลลัพธ์เหมือนกันทุกประการ แต่

เลข `100` อยู่ใน `domain` ที่เดียว ไม่มีใน Handler

`404`, `422`, `400` อยู่ใน `delivery/http/error.go` ที่เดียว ไม่มีใน Use Case

การเก็บข้อมูลอยู่ใน `repository` ที่เดียว ไม่มีใน Handler

</aside>

### 2.4.15 กฎการยกเลิก

ส่วนนี้คือเหตุผลทั้งหมดของการมีชั้น Use Case กฎมีข้อเดียวคือ *ยกเลิกได้ก่อนที่ร้านจะเริ่มทำอาหาร*

#### ขั้นที่ 20 — กฎการยกเลิกและการเดินสถานะ อยู่ในชั้น Use Case

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `usecase/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 7

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม Method `Cancel` ที่ถาม `order.CanCancel()` ก่อนตัดสินใจ
- เพิ่ม Method `SetStatus` ที่ถาม `order.CanAdvanceTo()` ก่อนเดินสถานะ

จากนั้นเพิ่มโค้ดด้านล่างนี้ต่อท้ายไฟล์

```go
// ตัวอย่างที่ชัดเจนที่สุดว่าเหตุใดจึงต้องมีชั้น Use Case
// เพราะไม่มี Constraint ใดในฐานข้อมูลที่เขียนกฎข้อนี้ได้
func (u OrderUseCase) Cancel(id string) error {
	order, err := u.repo.FindOrder(id)
	if err != nil {
		return fmt.Errorf("cancel order %s: %w", id, err)
	}

	if !order.CanCancel() {
		return domain.ErrNotCancelable
	}

	order.Status = domain.Cancelled
	return u.repo.Save(order)
}

func (u OrderUseCase) SetStatus(id string, next domain.Status) error {
	order, err := u.repo.FindOrder(id)
	if err != nil {
		return fmt.Errorf("set status %s: %w", id, err)
	}

	if !order.CanAdvanceTo(next) {
		return domain.ErrInvalidStatus
	}

	order.Status = next
	return u.repo.Save(order)
}
```

- 📄 ดู `usecase/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package usecase
    
    import (
    	"fmt"
    	"sync/atomic"
    
    	"food-api/domain"
    )
    
    // Interface ประกาศไว้ใน Package ที่เป็นผู้ใช้งาน ไม่ใช่ที่ผู้ให้บริการ
    // นี่คือ Dependency Inversion ทิศทางการพึ่งพาจึงชี้เข้าข้างใน
    type OrderRepository interface {
    	ListMenu() ([]domain.MenuItem, error)
    	FindMenuItem(id string) (domain.MenuItem, error)
    	FindOrder(id string) (domain.Order, error)
    	Save(o domain.Order) error
    }
    
    type OrderUseCase struct {
    	repo OrderRepository
    	seq  *int64
    }
    
    // รับ Interface เข้ามาเป็นพารามิเตอร์ แต่คืน Struct ตัวจริงออกไป
    func NewOrderUseCase(r OrderRepository) OrderUseCase {
    	var n int64
    	return OrderUseCase{repo: r, seq: &n}
    }
    
    func (u OrderUseCase) nextID() string {
    	return fmt.Sprintf("o%d", atomic.AddInt64(u.seq, 1))
    }
    
    func (u OrderUseCase) ListMenu() ([]domain.MenuItem, error) {
    	return u.repo.ListMenu()
    }
    
    // หนึ่ง Method ต่อหนึ่งสิ่งที่ผู้ใช้กระทำได้
    func (u OrderUseCase) Place(menuID string, qty int) (domain.Order, error) {
    	if qty < 1 {
    		return domain.Order{}, domain.ErrInvalidQty
    	}
    
    	item, err := u.repo.FindMenuItem(menuID)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    
    	order := domain.Order{
    		Lines:  []domain.OrderLine{{Item: item, Qty: qty}},
    		Status: domain.Pending,
    	}
    
    	// ตรวจสอบกฎให้ครบก่อน แล้วจึงกำหนด ID และบันทึก
    	// หากกำหนด ID ตั้งแต่ต้น ออเดอร์ที่ถูกปฏิเสธจะใช้หมายเลขสิ้นเปลืองโดยเปล่าประโยชน์
    	if order.Total() < domain.MinimumOrder {
    		return domain.Order{}, domain.ErrBelowMinimum
    	}
    
    	order.ID = u.nextID()
    
    	if err := u.repo.Save(order); err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    	return order, nil
    }
    
    func (u OrderUseCase) Get(id string) (domain.Order, error) {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("get order %s: %w", id, err)
    	}
    	return order, nil
    }
    
    // ตัวอย่างที่ชัดเจนที่สุดว่าเหตุใดจึงต้องมีชั้น Use Case
    // เพราะไม่มี Constraint ใดในฐานข้อมูลที่เขียนกฎข้อนี้ได้
    func (u OrderUseCase) Cancel(id string) error {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return fmt.Errorf("cancel order %s: %w", id, err)
    	}
    
    	if !order.CanCancel() {
    		return domain.ErrNotCancelable
    	}
    
    	order.Status = domain.Cancelled
    	return u.repo.Save(order)
    }
    
    func (u OrderUseCase) SetStatus(id string, next domain.Status) error {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return fmt.Errorf("set status %s: %w", id, err)
    	}
    
    	if !order.CanAdvanceTo(next) {
    		return domain.ErrInvalidStatus
    	}
    
    	order.Status = next
    	return u.repo.Save(order)
    }
    ```
    

#### ขั้นที่ 21 — DTO ของการเปลี่ยนสถานะ

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `delivery/http/dto.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 9

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม `setStatusRequest` สำหรับรับสถานะใหม่เข้ามา

ค้นหาส่วนนี้ในไฟล์ แล้วแก้ไขให้เป็นดังด้านล่าง โดยแสดงบรรทัดโดยรอบไว้เพื่อให้ระบุตำแหน่งได้ง่าย

```go
type placeOrderRequest struct {
	MenuID string `json:"menu_id"`
	Qty    int    `json:"qty"`
}

type setStatusRequest struct {
	Status string `json:"status"`
}
```

- 📄 ดู `delivery/http/dto.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/domain"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // รูปแบบของข้อมูลที่รับเข้ามา มี JSON Tag เพราะเป็นเรื่องของชั้นนอก
    type placeOrderRequest struct {
    	MenuID string `json:"menu_id"`
    	Qty    int    `json:"qty"`
    }
    
    type setStatusRequest struct {
    	Status string `json:"status"`
    }
    
    // รูปแบบของข้อมูลที่ส่งออกไป แยกจาก Type ใน domain
    // จึงเปลี่ยนชื่อ Key ใน JSON ได้โดยไม่ต้องแก้ไข Domain
    type orderResponse struct {
    	ID     string     `json:"id"`
    	Status string     `json:"status"`
    	Total  float64    `json:"total"`
    	Lines  []lineResp `json:"lines"`
    }
    
    type lineResp struct {
    	MenuID string  `json:"menu_id"`
    	Name   string  `json:"name"`
    	Price  float64 `json:"price"`
    	Qty    int     `json:"qty"`
    }
    
    type menuResponse struct {
    	ID    string  `json:"id"`
    	Name  string  `json:"name"`
    	Price float64 `json:"price"`
    }
    
    // Adapter ทำหน้าที่แปลงข้อมูลทางเดียว จาก domain ออกสู่ภายนอก
    func toResponse(o domain.Order) orderResponse {
    	lines := make([]lineResp, 0, len(o.Lines))
    	for _, l := range o.Lines {
    		lines = append(lines, lineResp{
    			MenuID: l.Item.ID,
    			Name:   l.Item.Name,
    			Price:  l.Item.Price,
    			Qty:    l.Qty,
    		})
    	}
    
    	return orderResponse{
    		ID:     o.ID,
    		Status: string(o.Status),
    		Total:  o.Total(),
    		Lines:  lines,
    	}
    }
    
    func toMenuResponse(items []domain.MenuItem) []menuResponse {
    	out := make([]menuResponse, 0, len(items))
    	for _, m := range items {
    		out = append(out, menuResponse{ID: m.ID, Name: m.Name, Price: m.Price})
    	}
    	return out
    }
    
    func msg(s string) fiber.Map { return fiber.Map{"error": s} }
    ```
    

#### ขั้นที่ 22 — Handler ของการยกเลิกและการเปลี่ยนสถานะ

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `delivery/http/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 10

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม `food-api/domain` เข้าไปในบล็อก `import` เพราะต้องแปลง string เป็น `domain.Status`
- เพิ่ม Handler `CancelOrder` และ `SetStatus`

แก้ไขบล็อก `import` ที่ส่วนหัวของไฟล์ ให้เป็นดังนี้

```go
import (
	"food-api/domain"
	"food-api/usecase"

	"github.com/gofiber/fiber/v2"
)
```

จากนั้นเพิ่มโค้ดด้านล่างนี้ต่อท้ายไฟล์

```go
func (h Handler) CancelOrder(c *fiber.Ctx) error {
	if err := h.uc.Cancel(c.Params("id")); err != nil {
		return toHTTP(c, err)
	}
	return c.SendStatus(204)
}

func (h Handler) SetStatus(c *fiber.Ctx) error {
	var req setStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(msg("invalid body"))
	}

	if err := h.uc.SetStatus(c.Params("id"), domain.Status(req.Status)); err != nil {
		return toHTTP(c, err)
	}
	return c.SendStatus(204)
}
```

- 📄 ดู `delivery/http/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/domain"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // Handler ถือ Use Case ไว้ ไม่ได้ถือการเชื่อมต่อฐานข้อมูล
    type Handler struct {
    	uc usecase.OrderUseCase
    }
    
    func NewHandler(uc usecase.OrderUseCase) Handler {
    	return Handler{uc: uc}
    }
    
    // ชั้นนี้ทำเพียงสามอย่าง คือแปลง Request เรียก Use Case หนึ่งตัว แล้วเขียน Response
    // ไม่ปรากฏเลข 100 อยู่ในไฟล์นี้แม้แต่แห่งเดียว
    func (h Handler) ListMenu(c *fiber.Ctx) error {
    	items, err := h.uc.ListMenu()
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toMenuResponse(items))
    }
    
    func (h Handler) PlaceOrder(c *fiber.Ctx) error {
    	var req placeOrderRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	order, err := h.uc.Place(req.MenuID, req.Qty)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	return c.Status(201).JSON(toResponse(order))
    }
    
    func (h Handler) GetOrder(c *fiber.Ctx) error {
    	order, err := h.uc.Get(c.Params("id"))
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toResponse(order))
    }
    
    func (h Handler) CancelOrder(c *fiber.Ctx) error {
    	if err := h.uc.Cancel(c.Params("id")); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    
    func (h Handler) SetStatus(c *fiber.Ctx) error {
    	var req setStatusRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	if err := h.uc.SetStatus(c.Params("id"), domain.Status(req.Status)); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    ```
    

#### ขั้นที่ 23 — เดินสถานะจนยกเลิกไม่ได้ แล้วสังเกตผลลัพธ์ 409

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 11

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม Route `PATCH /api/orders/:id/status` และ `POST /api/orders/:id/cancel`

ค้นหาส่วนนี้ในไฟล์ แล้วแก้ไขให้เป็นดังด้านล่าง โดยแสดงบรรทัดโดยรอบไว้เพื่อให้ระบุตำแหน่งได้ง่าย

```go
api.Get("/menu", h.ListMenu)
api.Post("/orders", h.PlaceOrder)
api.Get("/orders/:id", h.GetOrder)
api.Patch("/orders/:id/status", h.SetStatus)
api.Post("/orders/:id/cancel", h.CancelOrder)

log.Fatal(app.Listen(":3000"))
```

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"log"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    func seedMenu() map[string]domain.MenuItem {
    	return map[string]domain.MenuItem{
    		"m1": {ID: "m1", Name: "Pad Thai", Price: 80},
    		"m2": {ID: "m2", Name: "Tom Yum Goong", Price: 120},
    		"m3": {ID: "m3", Name: "Green Curry", Price: 95},
    		"m4": {ID: "m4", Name: "Thai Iced Tea", Price: 45},
    	}
    }
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึง memory, usecase และ http พร้อมกัน
    // ประกอบจากชั้นในออกมาสู่ชั้นนอก
    func main() {
    	// 1. ที่เก็บข้อมูล
    	repo := memory.NewOrderRepo(seedMenu())
    
    	// 2. กฎ
    	uc := usecase.NewOrderUseCase(repo)
    
    	// 3. ชั้นขนส่ง
    	h := http.NewHandler(uc)
    
    	app := fiber.New()
    
    	// Middleware ครอบทุก Request โดยเขียนไว้ที่เดียว (Decorator Pattern)
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	api := app.Group("/api")
    	api.Get("/menu", h.ListMenu)
    	api.Post("/orders", h.PlaceOrder)
    	api.Get("/orders/:id", h.GetOrder)
    	api.Patch("/orders/:id/status", h.SetStatus)
    	api.Post("/orders/:id/cancel", h.CancelOrder)
    
    	log.Fatal(app.Listen(":3000"))
    }
    ```
    

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -i -X PATCH localhost:3000/api/orders/o1/status -H 'Content-Type: application/json' -d '{"status":"cooking"}'
```

```bash
HTTP/1.1 409 Conflict
Content-Type: application/json

{"error":"invalid status transition"}
```

คำสั่งที่ 2

```bash
curl -i -X PATCH localhost:3000/api/orders/o1/status -H 'Content-Type: application/json' -d '{"status":"confirmed"}'
```

```bash
HTTP/1.1 204 No Content
```

คำสั่งที่ 3

```bash
curl -i -X PATCH localhost:3000/api/orders/o1/status -H 'Content-Type: application/json' -d '{"status":"cooking"}'
```

```bash
HTTP/1.1 204 No Content
```

คำสั่งที่ 4

```bash
curl -i -X POST localhost:3000/api/orders/o1/cancel
```

```bash
HTTP/1.1 409 Conflict
Content-Type: application/json

{"error":"cooking already started"}
```

<aside>
💡

สังเกตว่า `Cancel` ไม่ได้ตัดสินเองว่ายกเลิกได้หรือไม่ แต่ไปถาม `order.CanCancel()` ซึ่งอยู่ในชั้น `domain`

หากในอนาคตมีสถานะใหม่ชื่อ `packing` ที่ยกเลิกไม่ได้ จะแก้ไขที่ `CanCancel` **เพียงจุดเดียว** แล้วทุก Endpoint ที่เกี่ยวข้องจะถูกต้องตามไปด้วยทั้งหมด

</aside>

### 2.4.16 เพิ่ม Endpoint ใหม่ทั้งเส้น

ถึงตรงนี้ระบบทำงานครบแล้ว ส่วนนี้จะแสดงให้เห็นว่า *การเพิ่มความสามารถใหม่* ในโครงสร้างแบบนี้ ต้องแก้ไขกี่ไฟล์ และแก้ไขที่จุดใดบ้าง

โจทย์คือ `GET /api/orders/:id/total` ที่คืนยอดรวมของออเดอร์ รวมค่าจัดส่ง 30 บาท ทำ 3 ขั้น ชั้นละหนึ่งขั้น

#### ขั้นที่ 24 — เพิ่ม Endpoint ใหม่ ขั้นที่ 1 — Use Case

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `usecase/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 7

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม Method `Total` ที่รวมยอดแล้วบวกค่าจัดส่ง

จากนั้นเพิ่มโค้ดด้านล่างนี้ต่อท้ายไฟล์

```go
func (u OrderUseCase) Total(id string) (float64, error) {
	order, err := u.repo.FindOrder(id)
	if err != nil {
		return 0, fmt.Errorf("total %s: %w", id, err)
	}
	return order.Total() + domain.DeliveryFee, nil
}
```

- 📄 ดู `usecase/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package usecase
    
    import (
    	"fmt"
    	"sync/atomic"
    
    	"food-api/domain"
    )
    
    // Interface ประกาศไว้ใน Package ที่เป็นผู้ใช้งาน ไม่ใช่ที่ผู้ให้บริการ
    // นี่คือ Dependency Inversion ทิศทางการพึ่งพาจึงชี้เข้าข้างใน
    type OrderRepository interface {
    	ListMenu() ([]domain.MenuItem, error)
    	FindMenuItem(id string) (domain.MenuItem, error)
    	FindOrder(id string) (domain.Order, error)
    	Save(o domain.Order) error
    }
    
    type OrderUseCase struct {
    	repo OrderRepository
    	seq  *int64
    }
    
    // รับ Interface เข้ามาเป็นพารามิเตอร์ แต่คืน Struct ตัวจริงออกไป
    func NewOrderUseCase(r OrderRepository) OrderUseCase {
    	var n int64
    	return OrderUseCase{repo: r, seq: &n}
    }
    
    func (u OrderUseCase) nextID() string {
    	return fmt.Sprintf("o%d", atomic.AddInt64(u.seq, 1))
    }
    
    func (u OrderUseCase) ListMenu() ([]domain.MenuItem, error) {
    	return u.repo.ListMenu()
    }
    
    // หนึ่ง Method ต่อหนึ่งสิ่งที่ผู้ใช้กระทำได้
    func (u OrderUseCase) Place(menuID string, qty int) (domain.Order, error) {
    	if qty < 1 {
    		return domain.Order{}, domain.ErrInvalidQty
    	}
    
    	item, err := u.repo.FindMenuItem(menuID)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    
    	order := domain.Order{
    		Lines:  []domain.OrderLine{{Item: item, Qty: qty}},
    		Status: domain.Pending,
    	}
    
    	// ตรวจสอบกฎให้ครบก่อน แล้วจึงกำหนด ID และบันทึก
    	// หากกำหนด ID ตั้งแต่ต้น ออเดอร์ที่ถูกปฏิเสธจะใช้หมายเลขสิ้นเปลืองโดยเปล่าประโยชน์
    	if order.Total() < domain.MinimumOrder {
    		return domain.Order{}, domain.ErrBelowMinimum
    	}
    
    	order.ID = u.nextID()
    
    	if err := u.repo.Save(order); err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    	return order, nil
    }
    
    func (u OrderUseCase) Get(id string) (domain.Order, error) {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("get order %s: %w", id, err)
    	}
    	return order, nil
    }
    
    // ตัวอย่างที่ชัดเจนที่สุดว่าเหตุใดจึงต้องมีชั้น Use Case
    // เพราะไม่มี Constraint ใดในฐานข้อมูลที่เขียนกฎข้อนี้ได้
    func (u OrderUseCase) Cancel(id string) error {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return fmt.Errorf("cancel order %s: %w", id, err)
    	}
    
    	if !order.CanCancel() {
    		return domain.ErrNotCancelable
    	}
    
    	order.Status = domain.Cancelled
    	return u.repo.Save(order)
    }
    
    func (u OrderUseCase) SetStatus(id string, next domain.Status) error {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return fmt.Errorf("set status %s: %w", id, err)
    	}
    
    	if !order.CanAdvanceTo(next) {
    		return domain.ErrInvalidStatus
    	}
    
    	order.Status = next
    	return u.repo.Save(order)
    }
    
    func (u OrderUseCase) Total(id string) (float64, error) {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return 0, fmt.Errorf("total %s: %w", id, err)
    	}
    	return order.Total() + domain.DeliveryFee, nil
    }
    ```
    

#### ขั้นที่ 25 — เพิ่ม Endpoint ใหม่ ขั้นที่ 2 — Handler

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `delivery/http/order.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 10

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม Handler `OrderTotal`

จากนั้นเพิ่มโค้ดด้านล่างนี้ต่อท้ายไฟล์

```go
func (h Handler) OrderTotal(c *fiber.Ctx) error {
	id := c.Params("id")

	total, err := h.uc.Total(id)
	if err != nil {
		return toHTTP(c, err)
	}

	// ไม่มีการคำนวณใด ๆ ในชั้นนี้ การคำนวณทั้งหมดอยู่ใน domain และ usecase
	return c.JSON(fiber.Map{"id": id, "total": total})
}
```

- 📄 ดู `delivery/http/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/domain"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // Handler ถือ Use Case ไว้ ไม่ได้ถือการเชื่อมต่อฐานข้อมูล
    type Handler struct {
    	uc usecase.OrderUseCase
    }
    
    func NewHandler(uc usecase.OrderUseCase) Handler {
    	return Handler{uc: uc}
    }
    
    // ชั้นนี้ทำเพียงสามอย่าง คือแปลง Request เรียก Use Case หนึ่งตัว แล้วเขียน Response
    // ไม่ปรากฏเลข 100 อยู่ในไฟล์นี้แม้แต่แห่งเดียว
    func (h Handler) ListMenu(c *fiber.Ctx) error {
    	items, err := h.uc.ListMenu()
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toMenuResponse(items))
    }
    
    func (h Handler) PlaceOrder(c *fiber.Ctx) error {
    	var req placeOrderRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	order, err := h.uc.Place(req.MenuID, req.Qty)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	return c.Status(201).JSON(toResponse(order))
    }
    
    func (h Handler) GetOrder(c *fiber.Ctx) error {
    	order, err := h.uc.Get(c.Params("id"))
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toResponse(order))
    }
    
    func (h Handler) CancelOrder(c *fiber.Ctx) error {
    	if err := h.uc.Cancel(c.Params("id")); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    
    func (h Handler) SetStatus(c *fiber.Ctx) error {
    	var req setStatusRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	if err := h.uc.SetStatus(c.Params("id"), domain.Status(req.Status)); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    
    func (h Handler) OrderTotal(c *fiber.Ctx) error {
    	id := c.Params("id")
    
    	total, err := h.uc.Total(id)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	// ไม่มีการคำนวณใด ๆ ในชั้นนี้ การคำนวณทั้งหมดอยู่ใน domain และ usecase
    	return c.JSON(fiber.Map{"id": id, "total": total})
    }
    ```
    

#### ขั้นที่ 26 — เพิ่ม Endpoint ใหม่ ขั้นที่ 3 — Route

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 11

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- เพิ่ม Route `GET /api/orders/:id/total`

ค้นหาส่วนนี้ในไฟล์ แล้วแก้ไขให้เป็นดังด้านล่าง โดยแสดงบรรทัดโดยรอบไว้เพื่อให้ระบุตำแหน่งได้ง่าย

```go
api.Get("/menu", h.ListMenu)
api.Post("/orders", h.PlaceOrder)
api.Get("/orders/:id", h.GetOrder)
api.Get("/orders/:id/total", h.OrderTotal)
api.Patch("/orders/:id/status", h.SetStatus)
api.Post("/orders/:id/cancel", h.CancelOrder)
```

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"log"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    func seedMenu() map[string]domain.MenuItem {
    	return map[string]domain.MenuItem{
    		"m1": {ID: "m1", Name: "Pad Thai", Price: 80},
    		"m2": {ID: "m2", Name: "Tom Yum Goong", Price: 120},
    		"m3": {ID: "m3", Name: "Green Curry", Price: 95},
    		"m4": {ID: "m4", Name: "Thai Iced Tea", Price: 45},
    	}
    }
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึง memory, usecase และ http พร้อมกัน
    // ประกอบจากชั้นในออกมาสู่ชั้นนอก
    func main() {
    	// 1. ที่เก็บข้อมูล
    	repo := memory.NewOrderRepo(seedMenu())
    
    	// 2. กฎ
    	uc := usecase.NewOrderUseCase(repo)
    
    	// 3. ชั้นขนส่ง
    	h := http.NewHandler(uc)
    
    	app := fiber.New()
    
    	// Middleware ครอบทุก Request โดยเขียนไว้ที่เดียว (Decorator Pattern)
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	api := app.Group("/api")
    	api.Get("/menu", h.ListMenu)
    	api.Post("/orders", h.PlaceOrder)
    	api.Get("/orders/:id", h.GetOrder)
    	api.Get("/orders/:id/total", h.OrderTotal)
    	api.Patch("/orders/:id/status", h.SetStatus)
    	api.Post("/orders/:id/cancel", h.CancelOrder)
    
    	log.Fatal(app.Listen(":3000"))
    }
    ```
    

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

```bash
curl -i localhost:3000/api/orders/o1/total
```

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{"id":"o1","total":190}
```

<aside>
💡

สรุปการเพิ่ม Endpoint ใหม่หนึ่งเส้น แก้ไข 3 ไฟล์ ไฟล์ละไม่กี่บรรทัด และไม่ต้องแก้ไข `domain` กับ `repository` เลย

ที่สำคัญกว่านั้นคือ Handler ไม่มีการคำนวณแม้แต่แห่งเดียว ส่วน Use Case ไม่รู้จัก `fiber.Ctx` เลย

</aside>

### 2.4.17 เขียน Test ของกฎทางธุรกิจ

เป้าหมายของการแบ่งชั้นทั้งหมด สรุปรวบยอดอยู่ในสองไฟล์นี้

#### ขั้นที่ 27 — ทดสอบกฎทางธุรกิจ โดยไม่ต้องเปิดฐานข้อมูลและไม่ต้องเรียก HTTP

<aside>
📄

**1. สร้างไฟล์ที่ 1** `usecase/fakes_test.go`

**2. สร้างไฟล์ที่ 2** `usecase/order_test.go`

**3. รัน** `go test ./...`

ขั้นนี้มี **2 ไฟล์** สร้างให้ครบทุกไฟล์ก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

ไฟล์ที่ 1 — `usecase/fakes_test.go`

```go
package usecase_test

import "food-api/domain"

// ไฟล์นี้เก็บ "ตัวแทน" ของ Repository ที่ใช้เฉพาะตอนทดสอบ
// แยกออกมาจากไฟล์ที่เก็บ Test เพื่อให้ order_test.go เหลือแต่ตัวการทดสอบจริง
//
// ทั้งสอง Type ด้านล่างไม่ได้ประกาศว่า implements สิ่งใด
// เพียงมี Method ครบตามที่ usecase.OrderRepository ต้องการ ก็ใช้แทนกันได้ทันที
// คุณสมบัตินี้เรียกว่า Implicit Satisfaction

// stubRepo คืนค่าที่กำหนดไว้ล่วงหน้า และจดจำสิ่งที่ถูกบันทึกไว้ใน saved
type stubRepo struct {
	order domain.Order
	item  domain.MenuItem
	saved domain.Order
}

func (s stubRepo) ListMenu() ([]domain.MenuItem, error) {
	return []domain.MenuItem{s.item}, nil
}

func (s stubRepo) FindMenuItem(id string) (domain.MenuItem, error) {
	return s.item, nil
}

func (s stubRepo) FindOrder(id string) (domain.Order, error) {
	return s.order, nil
}

// ใช้ Pointer receiver เพราะต้องเขียนค่ากลับลงในฟีลด์ saved
func (s *stubRepo) Save(o domain.Order) error {
	s.saved = o
	return nil
}

// failingRepo ล้มเหลวเสมอ ใช้ทดสอบเส้นทางที่ผิดพลาด โดยไม่ต้องมีฐานข้อมูลจริง
type failingRepo struct{}

func (failingRepo) ListMenu() ([]domain.MenuItem, error) {
	return nil, domain.ErrNotFound
}

func (failingRepo) FindMenuItem(string) (domain.MenuItem, error) {
	return domain.MenuItem{}, domain.ErrNotFound
}

func (failingRepo) FindOrder(string) (domain.Order, error) {
	return domain.Order{}, domain.ErrNotFound
}

func (failingRepo) Save(domain.Order) error { return nil }
```

ไฟล์ที่ 2 — `usecase/order_test.go`

```go
package usecase_test

import (
	"errors"
	"testing"

	"food-api/domain"
	"food-api/usecase"
)

// ตัวแทนของ Repository ที่ใช้ในไฟล์นี้ ประกาศไว้ใน fakes_test.go
// Test ทั้งหมดด้านล่างทำงานได้โดยไม่ต้องมีฐานข้อมูล และไม่ต้องเปิดเซิร์ฟเวอร์

func TestCancelRefusedWhileCooking(t *testing.T) {
	repo := &stubRepo{order: domain.Order{ID: "o1", Status: domain.Cooking}}
	uc := usecase.NewOrderUseCase(repo)

	err := uc.Cancel("o1")

	if !errors.Is(err, domain.ErrNotCancelable) {
		t.Fatalf("ต้องการ ErrNotCancelable แต่ได้ %v", err)
	}
}

func TestCancelAllowedWhilePending(t *testing.T) {
	repo := &stubRepo{order: domain.Order{ID: "o1", Status: domain.Pending}}
	uc := usecase.NewOrderUseCase(repo)

	if err := uc.Cancel("o1"); err != nil {
		t.Fatalf("ไม่ควรเกิด error แต่ได้ %v", err)
	}
	if repo.saved.Status != domain.Cancelled {
		t.Errorf("สถานะที่บันทึกควรเป็น cancelled แต่ได้ %v", repo.saved.Status)
	}
}

func TestPlaceBelowMinimum(t *testing.T) {
	repo := &stubRepo{item: domain.MenuItem{ID: "m1", Name: "Water", Price: 20}}
	uc := usecase.NewOrderUseCase(repo)

	_, err := uc.Place("m1", 2) // 40 บาท ไม่ถึงยอดขั้นต่ำ 100 บาท

	if !errors.Is(err, domain.ErrBelowMinimum) {
		t.Fatalf("ต้องการ ErrBelowMinimum แต่ได้ %v", err)
	}
}

func TestPlaceRejectsZeroQty(t *testing.T) {
	repo := &stubRepo{item: domain.MenuItem{ID: "m1", Price: 500}}
	uc := usecase.NewOrderUseCase(repo)

	if _, err := uc.Place("m1", 0); !errors.Is(err, domain.ErrInvalidQty) {
		t.Fatalf("ต้องการ ErrInvalidQty แต่ได้ %v", err)
	}
}

func TestPlaceWrapsNotFound(t *testing.T) {
	uc := usecase.NewOrderUseCase(failingRepo{})

	_, err := uc.Place("m9", 2)

	// เมื่อห่อด้วย %w แล้ว errors.Is จะยังตรวจสอบทะลุถึงต้นเหตุได้
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("ต้องการ ErrNotFound แต่ได้ %v", err)
	}
	if err.Error() != "place order: not found" {
		t.Errorf("ข้อความควรมีบริบทเพิ่ม แต่ได้ %q", err.Error())
	}
}

func TestTotalAddsDeliveryFee(t *testing.T) {
	repo := &stubRepo{order: domain.Order{
		ID:     "o1",
		Status: domain.Pending,
		Lines: []domain.OrderLine{
			{Item: domain.MenuItem{Price: 80}, Qty: 2},
			{Item: domain.MenuItem{Price: 120}, Qty: 1},
		},
	}}
	uc := usecase.NewOrderUseCase(repo)

	got, err := uc.Total("o1")
	if err != nil {
		t.Fatal(err)
	}
	if got != 310 {
		t.Errorf("ต้องการ 310 แต่ได้ %v", got)
	}
}

// ล็อกพฤติกรรมไว้ ออเดอร์ที่ถูกปฏิเสธต้องไม่ใช้หมายเลขสิ้นเปลือง
func TestRejectedOrderDoesNotConsumeID(t *testing.T) {
	repo := &stubRepo{item: domain.MenuItem{ID: "m1", Name: "Water", Price: 20}}
	uc := usecase.NewOrderUseCase(repo)

	if _, err := uc.Place("m1", 1); !errors.Is(err, domain.ErrBelowMinimum) {
		t.Fatalf("ต้องการ ErrBelowMinimum แต่ได้ %v", err)
	}

	repo.item = domain.MenuItem{ID: "m2", Name: "Pad Thai", Price: 80}
	order, err := uc.Place("m2", 2)
	if err != nil {
		t.Fatal(err)
	}

	if order.ID != "o1" {
		t.Errorf("ออเดอร์แรกที่สำเร็จควรได้หมายเลข o1 แต่ได้ %q", order.ID)
	}
}
```

**ผลลัพธ์**

```bash
go test ./...
```

```bash
?   	food-api	[no test files]
?   	food-api/delivery/http	[no test files]
ok  	food-api/domain	(cached)
?   	food-api/repository/memory	[no test files]
ok  	food-api/usecase	(cached)
```

<aside>
💡

`stubRepo` ไม่ต้องประกาศว่า implements สิ่งใด เพียงมี Method ครบตามที่ `usecase` ต้องการก็ใช้แทนกันได้ นี่คือ Implicit Satisfaction ที่ทำให้การเขียน Test ในภาษา Go ทำได้ง่ายมาก

Test ทั้งชุดนี้ไม่ใช้ฐานข้อมูล ไม่เปิดเซิร์ฟเวอร์ ไม่เรียก `curl` และทำงานเสร็จในเวลาไม่ถึงหนึ่งวินาที

</aside>

### 2.4.18 Design Patterns in This Code

โค้ดที่เขียนมาทั้งหมดใช้ Design Pattern อยู่หลายตัวโดยไม่ได้เรียกชื่อ

| Pattern | อยู่ที่จุดใดในโค้ดของเรา |
| --- | --- |
| Repository | `FindMenuItem` และ `Save` ซ่อนที่เก็บข้อมูลไว้หลังสิ่งที่ดูเหมือนคอลเลกชัน |
| Dependency Injection | `NewOrderUseCase(repo)` รับสิ่งที่ต้องพึ่งพาเข้ามา แทนที่จะสร้างขึ้นเอง |
| Adapter | `postgres.OrderRepo` แปลง `pgx` ให้เข้ากับ Interface ของเรา ส่วน Handler แปลง HTTP ให้เป็นการเรียก Use Case |
| Strategy | `memory` หรือ `postgres` ใช้ Interface เดียวกัน เลือกตอนประกอบระบบ โดยไม่ต้องแก้โค้ดที่อื่น |
| DTO | `orderResponse` พาข้อมูลข้ามขอบเขต ทำให้ `domain.Order` ไม่หลุดรูปแบบออกไป |
| Decorator / Chain | `app.Use(...)` ใน `main.go` ครอบทุก Request ด้วย Logger ตัวเดียว |

อีกหนึ่ง Pattern ที่จะได้ใช้เมื่อระบบโตขึ้น

```go
// Functional Options — Go ไม่มี Default Parameter
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

srv := NewServer(WithTimeout(5 * time.Second))
```

<aside>
💡

Go ต้องใช้ Pattern น้อยกว่าภาษาเชิงวัตถุอื่น เพราะบางอย่างมีอยู่ในตัวภาษาแล้ว เช่น Iterator คือ `range`, Observer คือ Channel และ Template Method คือการส่งฟังก์ชันเข้าไปเป็น Parameter

</aside>

### 2.4.19 The Endpoints

สรุป Endpoint ทั้งหมดที่สร้างเสร็จแล้ว

| Method | Path | ทำหน้าที่ใด |
| --- | --- | --- |
| `GET` | `/api/menu` | ดูเมนู ส่งต่อไปยัง Repository ได้โดยตรง |
| `POST` | `/api/orders` | สั่งอาหาร กฎยอดขั้นต่ำและการตรวจเมนูทำงานที่นี่ |
| `GET` | `/api/orders/:id` | ดูออเดอร์หนึ่งรายการ หรือคืน `404` |
| `GET` | `/api/orders/:id/total` | ยอดรวมพร้อมค่าจัดส่ง |
| `PATCH` | `/api/orders/:id/status` | เดินสถานะไปข้างหน้า การข้ามสถานะที่ผิดกฎคืน `409` |
| `POST` | `/api/orders/:id/cancel` | ยกเลิก ยกเว้นเมื่อเริ่มทำอาหารแล้ว |

## 2.5 Data and Persistence

### 2.5.1 Stateless and Stateful

ความแตกต่างข้อนี้เป็นตัวกำหนดว่าระบบจะขยายและ Deploy ได้อย่างไร

| Stateless — `api` | Stateful — `db` |
| --- | --- |
| ไม่เก็บสถานะใดไว้ระหว่าง Request | ถือข้อมูลที่ทุกอย่างพึ่งพาอยู่ |
| สำเนาตัวใดก็รับ Request ได้เหมือนกัน | ตัวตนสำคัญ เพราะเป็นเจ้าของ Volume |
| ปิด เปิดใหม่ หรือเปลี่ยนสำเนาได้โดยอิสระ | เปลี่ยนสำเนาไม่ได้โดยอิสระ |
| ขยายด้วยการเพิ่มจำนวนสำเนา | การขยายระบบ หมายถึงต้องจัดการ Replica, Failover และ Backup |
| State อยู่ใน Request หรือในฐานข้อมูล | หากข้อมูลสูญหาย ธุรกิจจะหยุดชะงักทันที |

<aside>
⚠️

`food-api` ที่เขียนมาถึงขั้นที่ 27 ยังเป็น **Stateful** อยู่ เพราะออเดอร์เก็บอยู่ใน Map ในหน่วยความจำของ Process เมื่อกด `Ctrl` + `C` ออเดอร์จะสูญหายทั้งหมด

หัวข้อนี้จะย้าย State ออกไปไว้ที่ Postgres ทำให้ตัว API กลายเป็น Stateless จริง ๆ

</aside>

ในระบบจริง การแบ่งนี้ปรากฏออกมาดังนี้

1. **หลายสำเนาของ API** วางอยู่หลัง Load Balancer หากสำเนาใดหยุดทำงาน ก็ส่งต่อไปยังสำเนาอื่น
2. **ใช้ Postgres แบบ Managed** เช่น RDS, Cloud SQL หรือ Neon เพื่อให้คนอื่นดูแล Backup และ Failover
3. **เก็บไฟล์ไว้ที่ Object Storage** อย่าง S3 หรือ GCS ไม่ใช่ในไฟล์ระบบของ Container เพราะ Container ถูกทิ้งได้ตลอด
4. **เก็บ Session ไว้ใน JWT หรือ Redis** ไม่ใช่ในหน่วยความจำของ Process เพราะระบบจะทำงานผิดพลาดทันทีที่เพิ่มจำนวนสำเนา
5. **Cache ต้องสูญหายได้** หาก Redis ไม่พร้อมใช้งาน ระบบควรทำงานช้าลง ไม่ใช่หยุดทำงาน
6. **Kubernetes แยกสองอย่างนี้ชัดเจน** ใช้ Deployment กับงาน Stateless และใช้ StatefulSet คู่กับ Volume Claim กับงาน Stateful

### 2.5.2 Docker และ Docker Compose อย่างย่อ

ถึงตรงนี้ต้องมี Postgres สักตัวให้โปรแกรมต่อเข้าไป ซึ่งทำได้สองทาง

1. **ให้ Docker เปิดให้** ไม่ต้องติดตั้ง Postgres ลงเครื่อง ลบทิ้งได้ทั้งชุดเมื่อเลิกใช้ เป็นทางที่เอกสารนี้ใช้เป็นหลัก
2. **ใช้ Postgres ที่ติดตั้งไว้ในเครื่องอยู่แล้ว** เช่นที่ลงไว้จากวิชาฐานข้อมูล ก็ใช้ได้เช่นกัน

<aside>
💡

ผู้ที่มี Postgres ของตนเองอยู่แล้ว **ข้ามหัวข้อนี้และหัวข้อ 2.5.4 ไปได้เลย** ไม่ต้องติดตั้ง Docker

สิ่งที่ต้องทำมีเพียงแก้ค่า `POSTGRES_*` ในไฟล์ `.env` ที่จะสร้างในหัวข้อถัดไป ให้ตรงกับฐานข้อมูลของตนเอง

โค้ดของโปรเจกต์ไม่ต่างกันแม้แต่บรรทัดเดียว เพราะโปรแกรมรู้จักเพียง Connection String เส้นเดียวเท่านั้น

</aside>

คำศัพท์ที่ต้องแยกให้ออกมีเพียงสี่คำ

| คำ | ความหมาย |
| --- | --- |
| Image | แบบพิมพ์ที่อ่านอย่างเดียว เช่น `postgres:16-alpine` คือ Postgres รุ่น 16 ที่มีผู้ประกอบไว้ให้แล้ว ดาวน์โหลดครั้งเดียวใช้ได้ตลอด |
| Container | Image ที่ถูกเรียกให้ทำงานจริง เปรียบได้กับความสัมพันธ์ระหว่าง Struct กับ Instance ของ Struct นั้น Container ถูกลบแล้วสร้างใหม่ได้ตลอดเวลา |
| Volume | ที่เก็บข้อมูลซึ่งอยู่ **นอก** Container จึงมีอายุยืนกว่า Container ข้อมูลในฐานข้อมูลอยู่ที่นี่ ไม่ได้อยู่ใน Container |
| Port mapping | การเปิดช่องจากเครื่องของเราเข้าไปถึง Container เขียนเป็น `"5432:5432"` โดยเลขหน้าเครื่องหมาย `:` คือ Port บนเครื่องของเรา ส่วนเลขหลังคือ Port ข้างใน Container |

**Docker Compose** คือการเขียนสิ่งเหล่านี้ลงไฟล์ `docker-compose.yml` หนึ่งไฟล์ แทนการพิมพ์ `docker run` ที่มีตัวเลือกยาวหลายบรรทัดทุกครั้ง เมื่อเขียนไว้แล้วทั้งทีมจะได้ฐานข้อมูลที่ตั้งค่าเหมือนกันทุกคน และไฟล์นี้ Commit ลง Git ได้

คำสั่งที่ใช้จริงในวิชานี้มีเพียงเท่านี้

| คำสั่ง | ทำสิ่งใด |
| --- | --- |
| `docker compose up -d` | เปิด Service ทั้งหมดในไฟล์ โดย `-d` คือให้ทำงานอยู่เบื้องหลัง ไม่ยึด Terminal ไว้ |
| `docker compose ps` | ดูว่าตัวใดทำงานอยู่ และพร้อมรับการเชื่อมต่อแล้วหรือยัง |
| `docker compose logs db` | อ่าน Log ของ Service ชื่อ `db` ใช้เมื่อต่อไม่ติดแล้วต้องการทราบสาเหตุ |
| `docker compose exec db psql ...` | สั่งงานข้างใน Container ที่กำลังทำงานอยู่ ในที่นี้คือเปิด `psql` โดยไม่ต้องติดตั้ง `psql` ลงเครื่อง |
| `docker compose down` | ลบ Container ทิ้ง แต่ **ยังเก็บ Volume ไว้** ข้อมูลจึงไม่หาย |
| `docker compose down -v` | ลบ Volume ทิ้งด้วย **ข้อมูลหายถาวร** มีเพียงคำสั่งนี้คำสั่งเดียวที่ลบข้อมูล |

<aside>
🚨

ทุกคำสั่งด้านบนต้องรันในโฟลเดอร์ที่มีไฟล์ `docker-compose.yml` อยู่ ซึ่งก็คือโฟลเดอร์ `food-api`

หากรันจากโฟลเดอร์อื่น จะได้ข้อความว่า `no configuration file provided`

</aside>

<aside>
💡

วิธีติดตั้ง Docker Desktop ทั้งบน macOS และ Windows พร้อมวิธีตรวจสอบว่าติดตั้งสำเร็จแล้ว อยู่ในไฟล์ `INSTALL-DOCKER.md`

ก่อนรันคำสั่งใดในหัวข้อถัดไป ต้องเปิดโปรแกรม Docker Desktop ค้างไว้เสมอ ตรวจได้ด้วยคำสั่ง `docker info`

</aside>

### 2.5.3 ค่าตั้งต้นของโปรเจกต์ ในไฟล์ .env

ตั้งแต่หัวข้อถัดไป โปรแกรมต้องรู้จัก Connection String ของฐานข้อมูล ซึ่งเป็นค่าที่ **ต่างกันไปในแต่ละเครื่อง** และมีรหัสผ่านอยู่ข้างใน จึงเขียนฝังลงในโค้ดไม่ได้

วิธีที่มักพบใน Tutorial คือพิมพ์ค่าไว้หน้าคำสั่งทุกครั้ง แบบนี้

```bash
DATABASE_URL=postgres://app:secret@localhost:5432/food go run .
```

วิธีนั้นใช้อธิบายให้เห็นภาพได้ดี แต่ใช้ทำงานจริงไม่ไหว ด้วยเหตุผลสี่ข้อ

1. **รหัสผ่านตกไปอยู่ในประวัติคำสั่งของ Terminal** ซึ่งเป็นไฟล์ธรรมดาที่อ่านย้อนหลังได้
2. **ต้องพิมพ์ซ้ำทุกครั้ง** ทั้งตอนรัน Migration และตอนรันเซิร์ฟเวอร์ พิมพ์ผิดที่ใดที่หนึ่งก็ต่อไม่ติด
3. **ค่าต่างกันไปในแต่ละเครื่อง** เพื่อนที่ใช้ Postgres ของตนเองจะมีชื่อผู้ใช้ Port และชื่อฐานข้อมูลคนละชุดกับเรา
4. **Windows เขียนไม่เหมือนกัน** รูปแบบ `VAR=value คำสั่ง` เป็นของ Shell ตระกูล Unix ส่วน PowerShell ต้องเขียน `$env:VAR="value"` แยกอีกบรรทัด

วิธีมาตรฐานคือเก็บค่าเหล่านี้ไว้ในไฟล์ `.env` แล้วให้โปรแกรมอ่านเอง โดยแยกเป็นสามไฟล์ที่มีหน้าที่ต่างกัน

| ไฟล์ | ลง Git หรือไม่ | หน้าที่ |
| --- | --- | --- |
| `.env.example` | **ลง** | ไฟล์ตัวอย่าง บอกว่าโปรเจกต์นี้ต้องการค่าใดบ้าง ใส่ค่าสมมติไว้ ไม่มีรหัสผ่านจริง |
| `.env` | **ไม่ลง** | ค่าจริงของเครื่องเราคนเดียว คัดลอกมาจาก `.env.example` แล้วแก้ให้ตรงกับเครื่องของตน |
| `.gitignore` | **ลง** | บอก Git ว่าห้ามนำ `.env` ขึ้นไป เป็นตัวบังคับใช้กฎข้อบนจริง ๆ |

ไฟล์ `.env` ไฟล์เดียวนี้ถูกอ่านโดยสองฝ่าย โดยที่ทั้งสองฝ่ายไม่รู้จักกัน

- **Docker Compose** อ่านให้เองโดยอัตโนมัติ นำไปแทนค่าในตำแหน่งที่เขียน `${...}` ไว้ในไฟล์ Compose ไม่ต้องติดตั้งสิ่งใดเพิ่ม
- **โปรแกรม Go** อ่านผ่านไลบรารี `godotenv` ด้วยโค้ดบรรทัดเดียว คือ `godotenv.Load()`

<aside>
💡

แนวคิดนี้มาจากหลัก **Config in the environment** ของ Twelve-Factor App คือสิ่งที่เปลี่ยนไปตามที่ติดตั้ง ต้องอยู่นอกโค้ด ไม่ใช่ในโค้ด

ผลที่ได้คือโค้ดชุดเดียวกันทำงานได้ทั้งบนเครื่องของเรา บนเครื่องของเพื่อน และบน Server จริง โดยไม่ต้องแก้ไฟล์ใดเลย เปลี่ยนเพียงค่าใน `.env`

</aside>

**ขั้นที่ 1 ย้ายค่าเข้าไฟล์ `.env` ตรง ๆ ซึ่งยังไม่ DRY**

เมื่อย้ายค่าเข้าไฟล์แล้ว สิ่งที่ได้ในตอนแรกมักเป็นแบบนี้ คือเขียน Connection String เต็มเส้นลงไปตรง ๆ

```bash
POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=food
POSTGRES_PORT=5432

# ชื่อผู้ใช้ รหัสผ่าน Port และชื่อฐานข้อมูล ถูกเขียนซ้ำอีกรอบทั้งเส้น
DATABASE_URL=postgres://app:secret@localhost:5432/food
```

ไฟล์นี้ใช้งานได้จริง แต่ค่าเดียวกันถูกเขียนไว้ **สองแห่ง** เมื่อเปลี่ยน `POSTGRES_PORT` เป็น `5433` แล้วลืมแก้เลขในสตริงของ `DATABASE_URL` ระบบจะต่อฐานข้อมูลไม่ติด โดยไม่มีสิ่งใดเตือนล่วงหน้า

นี่คือข้อจำกัดของการใช้ `dotenv` เพียว ๆ ที่อ่านไฟล์ตามตัวอักษรอย่างเดียว จึงต้องพิมพ์ค่าที่ประกอบขึ้นแล้วลงไปเอง

**ขั้นที่ 2 ให้ไฟล์ `.env` อ้างถึงค่าของตัวเอง**

ไฟล์ `.env` อ้างถึงค่าที่ประกาศไว้ก่อนหน้าในไฟล์เดียวกันได้ ด้วยรูปแบบ `${...}` จึงประกอบ Connection String ขึ้นจากค่าที่มีอยู่แล้ว

```bash
POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=food
POSTGRES_PORT=5432
POSTGRES_HOST=localhost

DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}
```

<aside>
💡

ความสามารถนี้คือสิ่งที่ฝั่ง Node.js เรียกว่า **`dotenv-expand`** ซึ่งเป็นแพ็กเกจแยกที่ต้องติดตั้งเพิ่มจาก `dotenv` เพราะ `dotenv` เพียว ๆ ไม่แทนค่า `${...}` ให้

ในภาษา Go ไม่ต้องติดตั้งอะไรเพิ่ม เพราะ `godotenv` รองรับการแทนค่าให้อยู่แล้วตั้งแต่ต้น เขียน `${...}` ได้ทันทีจาก `godotenv.Load()` ตัวเดียว

</aside>

ผลคือเหลือ **แหล่งความจริงเพียงแห่งเดียว** การเปลี่ยน Port เป็น `5433` แก้เพียงบรรทัดเดียว แล้วทั้ง Docker Compose และโปรแกรม Go จะถูกต้องตามไปพร้อมกัน

ข้อควรระวังของการแทนค่าในไฟล์ `.env` มีสามข้อ

| ข้อควรระวัง | รายละเอียด |
| --- | --- |
| ต้องประกาศก่อนใช้ | `godotenv` อ่านไฟล์จากบนลงล่าง การอ้างถึงตัวแปรที่อยู่บรรทัดล่างกว่า จะได้ค่าว่างโดยไม่แจ้งเตือน จึงต้องวางบรรทัด `POSTGRES_*` ไว้เหนือ `DATABASE_URL` เสมอ |
| Single Quote ปิดการแทนค่า | `'${POSTGRES_DB}'` จะได้ข้อความตามตัวอักษร ส่วนแบบไม่ใส่ Quote และแบบ Double Quote จะถูกแทนค่าตามปกติ |
| แทนค่าจากไฟล์เดียวกันเท่านั้น | ไม่ได้อ่านค่าจาก Terminal ดังนั้น `POSTGRES_PORT=5433 go run .` จะ **ไม่** ทำให้ `DATABASE_URL` เปลี่ยนตาม หากต้องการสั่งจาก Terminal ให้กำหนด `DATABASE_URL` ทั้งเส้นไปเลย |

#### ขั้นที่ 28 — ค่าตั้งต้นของโปรเจกต์ ในไฟล์ .env

<aside>
📄

**1. สร้างไฟล์ที่ 1** `.env.example` แล้ววางโค้ดด้านล่างลงไป

**2. สร้างไฟล์ที่ 2** `.gitignore` แล้ววางโค้ดด้านล่างลงไป

**3. คัดลอกเป็นไฟล์ใช้งานจริง** `cp .env.example .env`

**4. ตรวจว่า `.env` ถูกกันไว้แล้วจริง** ด้วย `git status --short` ต้องไม่ปรากฏชื่อ `.env`

ขั้นนี้มี **2 ไฟล์** สร้างให้ครบทุกไฟล์ก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

ไฟล์ที่ 1 — `.env.example`

```bash
# ไฟล์ตัวอย่าง Commit ลง Git ได้ เพราะไม่มีค่าจริงอยู่ข้างใน
# ก่อนเริ่มงาน ให้คัดลอกไฟล์นี้เป็น .env ก่อน   cp .env.example .env

# ─── แหล่งความจริงแห่งเดียว แก้ที่นี่ที่เดียวพอ ───
# ผู้ที่ใช้ Docker: Docker Compose อ่านสี่ค่าแรกไปสร้าง Container ให้เอง
# ผู้ที่มี Postgres อยู่แล้ว: แก้ห้าบรรทัดนี้ให้ตรงกับเครื่องของตน แล้วข้าม Docker ไปได้เลย
POSTGRES_USER=app
POSTGRES_PASSWORD=secret
POSTGRES_DB=food
POSTGRES_PORT=5432
POSTGRES_HOST=localhost

# ─── ค่าที่ประกอบขึ้นจากห้าบรรทัดด้านบน ───
# ${...} คือการอ้างถึงค่าที่ประกาศไว้ก่อนหน้าในไฟล์เดียวกัน
# godotenv แทนค่าให้เองตอนอ่านไฟล์ จึงไม่ต้องเขียนชื่อผู้ใช้และรหัสผ่านซ้ำอีกรอบ
# เมื่อแก้ค่าด้านบน บรรทัดนี้จะถูกต้องตามไปเอง ไม่ต้องแก้ซ้ำสองที่
DATABASE_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}

# หากลบบรรทัด DATABASE_URL ทิ้ง โปรแกรมจะกลับไปเก็บข้อมูลไว้ใน Memory
```

ไฟล์ที่ 2 — `.gitignore`

```bash
# ห้าม Commit ค่าจริง โดยเฉพาะรหัสผ่านและ Connection String
# ให้ Commit เฉพาะ .env.example ที่เป็นไฟล์ตัวอย่าง
.env

# ไฟล์ที่ได้จากการ Build
/food-api
/migrate
```

- ค่าที่พิมพ์นำหน้าคำสั่งใน Terminal **ชนะ** ค่าในไฟล์ `.env` เสมอ เพราะ `godotenv.Load()` จะไม่เขียนทับตัวแปรที่มีอยู่แล้ว จึงยังใช้สั่งชั่วคราวเพื่อทดลองได้
- `.env.example` ต้องอัปเดตทุกครั้งที่เพิ่มค่าตั้งต้นตัวใหม่ มิฉะนั้นเพื่อนที่ Clone โปรเจกต์ไปจะไม่ทราบว่าต้องกำหนดค่าใดบ้าง
- หากลืมสร้าง `.env` โปรแกรมจะไม่ล้มเหลวทันที แต่จะกลับไปเก็บข้อมูลใน Memory เพราะไม่พบ `DATABASE_URL` ให้สังเกตบรรทัดแรกใน Terminal เสมอ
- `DATABASE_URL` ประกอบขึ้นจากบรรทัด `POSTGRES_*` ด้านบนด้วย `${...}` จึงแก้ค่าที่เดียว ไม่ต้องแก้ซ้ำในสตริงของ Connection String อีกรอบ

<aside>
🚨

`.gitignore` ป้องกันได้เฉพาะไฟล์ที่ **ยังไม่เคยถูก Commit** เท่านั้น

หาก `.env` เคยถูก Commit ขึ้นไปแล้ว การเพิ่มชื่อลง `.gitignore` ทีหลังจะไม่มีผล ต้องสั่ง `git rm --cached .env` ก่อน

และเมื่อรหัสผ่านหลุดขึ้น Git ไปแล้ว ต้องถือว่ารหัสนั้นใช้ไม่ได้อีกต่อไป ให้เปลี่ยนรหัสใหม่เสมอ

</aside>

<aside>
💡

**ผู้ที่มี Postgres ของตนเองอยู่แล้ว** ให้แก้ค่า `POSTGRES_*` ทั้งห้าบรรทัดใน `.env` ให้ตรงกับเครื่องของตน แล้ว **ข้ามหัวข้อ 2.5.4 ไปที่ 2.5.5 ได้เลย**

ไม่ต้องแก้บรรทัด `DATABASE_URL` เลย เพราะประกอบขึ้นจากห้าค่านั้นให้เองอยู่แล้ว

ก่อนข้ามไป ต้องสร้างฐานข้อมูลเปล่าไว้หนึ่งชุดก่อน เพราะไฟล์ Migration สร้างเฉพาะตารางเท่านั้น ไม่ได้สร้างตัวฐานข้อมูล

`createdb food` หรือสั่ง `CREATE DATABASE food;` ใน `psql` ก็ได้เช่นกัน

</aside>

### 2.5.4 Postgres with Docker Compose

#### ขั้นที่ 29 — Postgres ด้วย Docker Compose

<aside>
📄

**1. สร้างไฟล์** `docker-compose.yml` แล้ววางโค้ดด้านล่างลงไป

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```yaml
# ไฟล์นี้ใช้เฉพาะผู้ที่ต้องการให้ Docker เปิด Postgres ให้
# ผู้ที่มี Postgres ติดตั้งไว้ในเครื่องอยู่แล้ว ไม่ต้องใช้ไฟล์นี้เลย
# เพียงแก้ค่า POSTGRES_* ใน .env ให้ตรงกับเครื่องของตนเอง แล้ว DATABASE_URL จะถูกต้องตามไปเอง

services:
  db:
    image: postgres:16-alpine
    environment:
      # อ่านมาจากไฟล์ .env ที่อยู่ข้าง ๆ ไฟล์นี้
      # Docker Compose อ่าน .env ให้เองโดยไม่ต้องติดตั้งไลบรารีใด
      # ค่าหลัง :- คือค่าที่ใช้เมื่อไม่พบตัวแปรนั้นในไฟล์ .env
      POSTGRES_USER: ${POSTGRES_USER:-app}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-secret}
      POSTGRES_DB: ${POSTGRES_DB:-food}
    ports:
      # รูปแบบคือ  "Port บนเครื่องเรา : Port ข้างใน Container"
      #
      #   เลขหน้าเครื่องหมาย :  คือ Port บนเครื่องของเรา เป็นเลขที่ใช้ต่อจากภายนอก
      #                        และเป็นเลขที่ DATABASE_URL ใน .env ดึงไปใช้ผ่าน ${POSTGRES_PORT} ด้วย
      #   เลขหลังเครื่องหมาย : คือ Port ข้างใน Container ซึ่ง Postgres ฟังอยู่ที่ 5432 เสมอ
      #                        เลขนี้เป็นของ Image ไม่ใช่ของเรา จึงห้ามแก้
      #
      # เมื่อ Port 5432 บนเครื่องถูกใช้อยู่แล้ว ให้แก้ POSTGRES_PORT ใน .env เป็น 5433
      # แก้บรรทัดเดียวพอ เพราะ DATABASE_URL ประกอบขึ้นจากค่านั้นอยู่แล้ว
      # ส่วนเลข 5432 หลังเครื่องหมาย : คงไว้เหมือนเดิม
      - "${POSTGRES_PORT:-5432}:5432"
    volumes:
      # Named Volume มีอายุยืนกว่า Container ข้อมูลจึงอยู่ข้ามการปิดเปิด
      - pgdata:/var/lib/postgresql/data
      #
      # สังเกตว่าไม่ได้ Mount โฟลเดอร์ migrations เข้า docker-entrypoint-initdb.d
      # เพราะสคริปต์ในโฟลเดอร์นั้นทำงานเฉพาะตอน Volume ว่างเปล่าเท่านั้น
      # จึงใช้เป็นเครื่องมือ Migration ไม่ได้ ให้ใช้ ./cmd/migrate แทน
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-app} -d ${POSTGRES_DB:-food}"]
      interval: 5s
      timeout: 3s
      retries: 10

volumes:
  pgdata:
```

- ใช้เพียงสองคำสั่งคือ `docker compose up -d` เพื่อเริ่ม และ `docker compose down` เพื่อหยุด ไม่ต้องติดตั้ง Postgres ลงเครื่อง
- ทุกตำแหน่งที่เขียน `${...}` ถูกแทนค่าจากไฟล์ `.env` ที่สร้างไว้ในขั้นที่ 28 รหัสผ่านจึงไม่ปรากฏในไฟล์ที่ Commit ลง Git
- ค่าหลัง `:-` คือค่าที่ใช้เมื่อไม่พบตัวแปรนั้น เช่น `${POSTGRES_PORT:-5432}` แปลว่า *ถ้าไม่กำหนด `POSTGRES_PORT` ให้ใช้ `5432`*
- `"${POSTGRES_PORT:-5432}:5432"` คือ Port mapping เลข **หน้า** เครื่องหมาย `:` คือ Port บนเครื่องของเรา ซึ่ง `DATABASE_URL` ใน `.env` ดึงไปใช้ผ่าน `${POSTGRES_PORT}` อยู่แล้ว จึงแก้ที่ `.env` ที่เดียว ส่วนเลข **หลัง** คือ Port ข้างใน Container ที่ Postgres ฟังอยู่เสมอ จึงไม่ต้องแก้
- ไม่ได้ Mount โฟลเดอร์ `migrations` เข้า `docker-entrypoint-initdb.d` เพราะสคริปต์ในโฟลเดอร์นั้นทำงานเฉพาะตอนที่ Volume ว่างเปล่า จึงใช้เป็นเครื่องมือ Migration ไม่ได้ เหตุผลเต็มอยู่ในหัวข้อ 2.5.7
- `healthcheck` ทำให้ `docker compose ps` บอกได้ว่าฐานข้อมูลพร้อมรับการเชื่อมต่อจริงหรือยัง ไม่ใช่เพียงการที่ Container เริ่มทำงานแล้วเท่านั้น

<aside>
💡

การใช้ `depends_on` เพียงอย่างเดียว จะรอเพียงให้ Container เริ่มทำงาน ไม่ได้รอให้ Postgres พร้อมรับการเชื่อมต่อ จึงต้องใช้คู่กับ `condition: service_healthy` เมื่อมี Service อื่นที่ต้องรอ

</aside>

### 2.5.5 Schema and Migrations

#### ขั้นที่ 30 — Schema และข้อมูลเริ่มต้น

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p migrations`

**2. สร้างไฟล์** `migrations/001_init.sql` แล้ววางโค้ดด้านล่างลงไป

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```sql
-- Migration ไฟล์แรก สร้างตารางทั้งหมด แล้วใส่ข้อมูลเมนูเริ่มต้น
-- ถูกรันโดยตัวรัน Migration ที่ ./cmd/migrate ไม่ใช่โดย Postgres เอง

CREATE TABLE menu_items (
    id    TEXT PRIMARY KEY,
    name  TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL
);

CREATE TABLE orders (
    id      TEXT PRIMARY KEY,
    status  TEXT NOT NULL,
    total   NUMERIC(10,2) NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_lines (
    order_id TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_id  TEXT NOT NULL REFERENCES menu_items(id),
    qty      INT  NOT NULL CHECK (qty > 0)
);

CREATE INDEX idx_order_lines_order_id ON order_lines(order_id);

-- ข้อมูลเมนูชุดเดียวกับที่ seedMenu() ใน main.go ใช้
INSERT INTO menu_items (id, name, price) VALUES
    ('m1', 'Pad Thai',      80),
    ('m2', 'Tom Yum Goong', 120),
    ('m3', 'Green Curry',   95),
    ('m4', 'Thai Iced Tea', 45);
```

```mermaid
erDiagram
    menu_items {
        TEXT id PK
        TEXT name
        NUMERIC price
    }
    orders {
        TEXT id PK
        TEXT status
        NUMERIC total
        TIMESTAMPTZ created
    }
    order_lines {
        TEXT order_id FK
        TEXT menu_id FK
        INT qty
    }
    orders ||--o{ order_lines : "มีหลายรายการ"
    menu_items ||--o{ order_lines : "ถูกอ้างถึงโดย"
```

- เงินใช้ `NUMERIC` ไม่ใช่ `float64` เพราะทศนิยมฐานสองปัดเศษไม่ตรงกับที่ผู้ใช้คาดหวัง
- `CHECK (qty > 0)` เป็นแนวป้องกันชั้นที่สอง ไม่ใช่ที่อยู่ของกฎ ตัวกฎจริงอยู่ใน `usecase.Place`

<aside>
🚨

ไฟล์นี้ **ยังไม่ถูกรัน** ในตอนนี้ เพราะไฟล์ Compose ของเราไม่ได้ Mount โฟลเดอร์ `migrations` เข้า `docker-entrypoint-initdb.d`

เหตุผลอยู่ในหัวข้อ 2.5.7 ถัดไป

</aside>

### 2.5.6 เปิดฐานข้อมูล

#### ขั้นที่ 31 — เปิดฐานข้อมูล แล้วสังเกตว่ายังไม่มีตารางใด ๆ

<aside>
📄

**1. เปิด Docker Desktop ให้เรียบร้อยก่อน** ตรวจด้วย `docker info`

**2. รัน** `docker compose up -d`

**3. รอจนกว่าสถานะจะเป็น `(healthy)`** ตรวจด้วย `docker compose ps`

**4. ตรวจดูตารางในฐานข้อมูล** ด้วยคำสั่ง `psql` ด้านล่าง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ผลลัพธ์**

คำสั่งที่ 1

```bash
docker compose up -d
```

```bash
Network food-api_default  Created
Volume food-api_pgdata    Created
Container food-api-db-1   Started
```

คำสั่งที่ 2

```bash
docker compose ps
```

```bash
SERVICE   STATUS
db        Up 5 seconds (healthy)
```

คำสั่งที่ 3

```bash
docker compose exec db psql -U app -d food -c "\dt"
```

```bash
Did not find any relations.
```

<aside>
💡

สังเกตว่าคำสั่ง `\dt` ตอบว่า `Did not find any relations.` แปลว่ายังไม่มีตารางใดถูกสร้างขึ้น

ทั้งที่เราเขียนไฟล์ `001_init.sql` ไปแล้วในขั้นที่ 30

</aside>

### 2.5.7 Migration คืออะไร และเหตุใดจึงต้องมี

เอกสารทั่วไปมักสอนให้ Mount โฟลเดอร์ `migrations` เข้า `docker-entrypoint-initdb.d` แล้วปล่อยให้ Postgres รันไฟล์ `.sql` ให้เอง

วิธีนั้นใช้ได้เฉพาะ **ครั้งแรกที่ Volume ว่างเปล่า** เท่านั้น

| สถานการณ์ | `docker-entrypoint-initdb.d` ทำงานหรือไม่ |
| --- | --- |
| เปิดฐานข้อมูลครั้งแรก Volume ว่างเปล่า | ทำงาน รันไฟล์ `.sql` ทั้งหมดหนึ่งครั้ง |
| ระบบใช้งานมาสักพัก แล้วเพิ่มไฟล์ `.sql` ใหม่ | **ไม่ทำงาน** ไฟล์ใหม่จะถูกข้ามไปเงียบ ๆ |
| Deploy ขึ้น Server ที่มีข้อมูลลูกค้าอยู่แล้ว | **ไม่ทำงาน** และเป็นสถานการณ์ที่พบบ่อยที่สุดในงานจริง |

งานจริงจึงต้องมี **ตัวรัน Migration** ที่ทำสองอย่าง

1. **จดจำ** ว่าไฟล์ใดถูกรันไปแล้ว โดยบันทึกไว้ในตารางของฐานข้อมูลเอง
2. **รันเฉพาะไฟล์ที่ยังไม่เคยรัน** เรียงตามลำดับเลขนำหน้า

เครื่องมือสำเร็จรูปอย่าง `golang-migrate`, `goose` และ `atlas` ทำงานด้วยหลักการนี้ทั้งหมด

ในวิชานี้จะเขียนขึ้นเองประมาณ 100 บรรทัด เพื่อให้เห็นว่าเบื้องหลังไม่มีอะไรซับซ้อน และไม่ต้องติดตั้งเครื่องมือเพิ่ม

#### ขั้นที่ 32 — เขียนตัวรัน Migration ขึ้นมาเอง

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p cmd/migrate`

**2. ติดตั้ง Library** `go get github.com/jackc/pgx/v5 github.com/joho/godotenv`

**3. สร้างไฟล์** `cmd/migrate/main.go` แล้ววางโค้ดด้านล่างลงไป

**4. รัน** `go run ./cmd/migrate`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
// ตัวรัน Migration ขนาดเล็ก เขียนเองเพื่อให้เห็นว่าเครื่องมือจริงทำงานอย่างไร
//
// รัน: go run ./cmd/migrate
// ต้องรันจากโฟลเดอร์ food-api เพราะอ่านทั้งไฟล์ .env และโฟลเดอร์ migrations จากตำแหน่งที่รัน

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const dir = "migrations"

func main() {
	// อ่านค่าจากไฟล์ .env เช่นเดียวกับ main.go ของตัวเซิร์ฟเวอร์
	// ไม่ตรวจสอบ Error เพราะผู้ที่กำหนดค่าจาก Terminal เองไม่จำเป็นต้องมีไฟล์นี้
	_ = godotenv.Load()

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("ไม่พบ DATABASE_URL ให้คัดลอก .env.example เป็น .env แล้วแก้ค่าให้ตรงกับฐานข้อมูลของตนเอง")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Fatal("เชื่อมต่อฐานข้อมูลไม่สำเร็จ: ", err)
	}
	defer conn.Close(ctx)

	// ตารางนี้คือสมุดบันทึกว่า Migration ใดถูกรันไปแล้วบ้าง
	// เครื่องมือจริงอย่าง golang-migrate หรือ goose ก็ใช้หลักการเดียวกัน
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	if err != nil {
		log.Fatal(err)
	}

	applied, err := appliedVersions(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	files, err := sqlFiles()
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for _, f := range files {
		version := strings.TrimSuffix(filepath.Base(f), ".sql")
		if applied[version] {
			continue
		}

		body, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}

		// หนึ่ง Migration หนึ่ง Transaction
		// หากไฟล์ใดล้มเหลวกลางคัน ฐานข้อมูลจะย้อนกลับทั้งไฟล์ ไม่ค้างครึ่ง ๆ กลาง ๆ
		if err := applyOne(ctx, conn, version, string(body)); err != nil {
			log.Fatalf("migration %s ล้มเหลว: %v", version, err)
		}

		log.Printf("ใช้งาน %s แล้ว", version)
		count++
	}

	if count == 0 {
		log.Println("ไม่มี migration ใหม่ ฐานข้อมูลเป็นรุ่นล่าสุดอยู่แล้ว")
	} else {
		log.Printf("ใช้งาน migration ใหม่ทั้งหมด %d ไฟล์", count)
	}
}

func appliedVersions(ctx context.Context, conn *pgx.Conn) (map[string]bool, error) {
	rows, err := conn.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

// เรียงตามชื่อไฟล์ เลขนำหน้าจึงเป็นตัวกำหนดลำดับการรัน
func sqlFiles() ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(out)
	return out, nil
}

func applyOne(ctx context.Context, conn *pgx.Conn, version, body string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, body); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
```

**ผลลัพธ์**

คำสั่งที่ 1

```bash
go run ./cmd/migrate
```

```bash
ใช้งาน 001_init แล้ว
ใช้งาน migration ใหม่ทั้งหมด 1 ไฟล์
```

คำสั่งที่ 2

```bash
docker compose exec db psql -U app -d food -c "\dt"
```

```bash
             List of relations
 Schema |       Name        | Type  | Owner
--------+-------------------+-------+-------
 public | menu_items        | table | app
 public | order_lines       | table | app
 public | orders            | table | app
 public | schema_migrations | table | app
(4 rows)
```

<aside>
💡

สังเกตว่าตอนนี้มีตาราง `schema_migrations` เพิ่มขึ้นมา ซึ่งไม่ได้เขียนไว้ในไฟล์ `.sql` ใด ๆ

ตัวรัน Migration สร้างขึ้นเองเพื่อใช้เป็นสมุดบันทึกว่าไฟล์ใดถูกรันไปแล้วบ้าง

</aside>

### 2.5.8 แก้ Schema บนฐานข้อมูลที่มีข้อมูลอยู่แล้ว

ส่วนนี้คือหัวใจของเรื่อง Migration ทั้งหมด สถานการณ์คือ **ระบบเปิดใช้งานไปแล้ว มีออเดอร์ของลูกค้าอยู่ในฐานข้อมูล**

แล้วร้านค้าขอให้เพิ่มช่องข้อความถึงร้าน

ในสถานการณ์แบบนี้จะลบฐานข้อมูลทิ้งแล้วสร้างใหม่ไม่ได้ ต้องแก้ Schema โดยที่ข้อมูลเดิมยังอยู่ครบ

#### ขั้นที่ 33 — Migration ที่ 2 — เพิ่มคอลัมน์บนตารางที่มีข้อมูลอยู่แล้ว

<aside>
📄

**1. สร้างไฟล์** `migrations/002_add_note_to_orders.sql` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go run ./cmd/migrate`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```sql
-- เพิ่มช่องข้อความถึงร้าน เช่น "ไม่ใส่ผักชี"
--
-- คอลัมน์ใหม่ต้องเป็น Nullable เสมอ เพราะแถวที่มีอยู่แล้วไม่มีค่านี้
-- หากใส่ NOT NULL โดยไม่มี DEFAULT คำสั่งจะล้มเหลวทันทีบนตารางที่มีข้อมูล
ALTER TABLE orders ADD COLUMN note TEXT;

-- ช่องเก็บข้อมูลที่รูปแบบยังไม่นิ่ง เช่นหมายเลขโต๊ะ หรือ Coupon Code
ALTER TABLE orders ADD COLUMN attributes JSONB NOT NULL DEFAULT '{}';
```

**ผลลัพธ์**

```bash
go run ./cmd/migrate
```

```bash
ใช้งาน 002_add_note_to_orders แล้ว
ใช้งาน migration ใหม่ทั้งหมด 1 ไฟล์
```

<aside>
💡

สังเกตว่าตัวรัน Migration ใช้งานเฉพาะไฟล์ `002` ไฟล์เดียว ไม่ได้รัน `001` ซ้ำอีกครั้ง

เพราะ `001` ถูกบันทึกไว้ใน `schema_migrations` แล้ว ข้อมูลเดิมในตาราง `orders` จึงไม่ถูกแตะต้อง

นี่คือความหมายของกฎ **เพิ่มได้ แต่อย่าลบ** ในหัวข้อ 2.5.12

</aside>

<aside>
🚨

คอลัมน์ `note` ประกาศเป็น Nullable และ `attributes` มี `DEFAULT '{}'`

หากเขียนเป็น `NOT NULL` โดยไม่มี `DEFAULT` คำสั่งจะล้มเหลวทันที เพราะแถวที่มีอยู่แล้วไม่มีค่าให้ใส่

`ERROR: column "note" of relation "orders" contains null values`

</aside>

#### ขั้นที่ 34 — Migration ที่ 3 — เพิ่ม Index แล้วตรวจสมุดบันทึก

<aside>
📄

**1. สร้างไฟล์** `migrations/003_index_orders_status.sql` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go run ./cmd/migrate`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```sql
-- หน้าจอของร้านต้องค้นออเดอร์ตามสถานะบ่อยมาก เช่น "ขอรายการที่ยัง pending ทั้งหมด"
-- หากไม่มี Index ฐานข้อมูลต้องไล่อ่านทุกแถวในตาราง
CREATE INDEX idx_orders_status ON orders(status);

-- Index สำหรับเรียงตามเวลาที่สั่ง ใช้กับหน้าจอที่แสดงออเดอร์ล่าสุดขึ้นก่อน
CREATE INDEX idx_orders_created ON orders(created DESC);
```

**ผลลัพธ์**

คำสั่งที่ 1

```bash
go run ./cmd/migrate
```

```bash
ใช้งาน 003_index_orders_status แล้ว
ใช้งาน migration ใหม่ทั้งหมด 1 ไฟล์
```

คำสั่งที่ 2

```bash
go run ./cmd/migrate
```

```bash
ไม่มี migration ใหม่ ฐานข้อมูลเป็นรุ่นล่าสุดอยู่แล้ว
```

คำสั่งที่ 3

```bash
docker compose exec db psql -U app -d food -c "SELECT version, applied_at FROM schema_migrations ORDER BY version;"
```

```bash
         version         |          applied_at
-------------------------+-------------------------------
 001_init                | 2026-09-03 08:42:09.126725+00
 002_add_note_to_orders  | 2026-09-03 08:42:09.314608+00
 003_index_orders_status | 2026-09-03 08:42:09.403289+00
(3 rows)
```

<aside>
💡

การรันครั้งที่สองตอบว่า *ไม่มี migration ใหม่* แปลว่าคำสั่งนี้ปลอดภัยที่จะรันซ้ำกี่ครั้งก็ได้

คุณสมบัตินี้เรียกว่า Idempotent และจำเป็นมากเวลา Deploy อัตโนมัติ เพราะสคริปต์ Deploy จะเรียกคำสั่งนี้ทุกครั้งโดยไม่ต้องตรวจก่อน

</aside>

<aside>
🧠

**เจาะลึก** แต่ละไฟล์ถูกรันอยู่ใน Transaction ของตัวเอง

หากไฟล์ใดมีคำสั่งผิดอยู่กลางไฟล์ ฐานข้อมูลจะย้อนกลับทั้งไฟล์ และจะ **ไม่บันทึก** ลง `schema_migrations`

ผลคือ Schema ไม่ค้างอยู่ในสภาพครึ่ง ๆ กลาง ๆ และเมื่อแก้ไฟล์แล้วรันใหม่ ระบบจะลองไฟล์นั้นอีกครั้งให้เอง

ทดลองได้โดยสร้างไฟล์ `004_broken.sql` ที่มีคำสั่งอ้างถึงตารางที่ไม่มีอยู่จริง แล้วรันดู จากนั้นลบไฟล์นั้นทิ้ง

</aside>

**กฎที่ต้องจำเกี่ยวกับไฟล์ Migration**

1. **ห้ามแก้ไฟล์ที่รันไปแล้ว** เพราะเครื่องอื่นที่รันไปแล้วจะไม่รันซ้ำ ทำให้ Schema ไม่ตรงกัน หากต้องเปลี่ยน ให้เพิ่มไฟล์ใหม่เสมอ
2. **เลขนำหน้าต้องเรียงและไม่ซ้ำ** เพราะลำดับการรันมาจากการเรียงชื่อไฟล์
3. **หนึ่งไฟล์ควรทำเรื่องเดียว** เพื่อให้ย้อนดูภายหลังได้ว่าการเปลี่ยนแปลงใดเกิดขึ้นเมื่อใด
4. **Commit ไฟล์ Migration ลง Git เสมอ** เพราะเป็นส่วนหนึ่งของโค้ด ไม่ใช่ของชั่วคราว

### 2.5.9 The Swap

ขั้นตอนที่พิสูจน์ว่าการแบ่งชั้นทั้งหมดคุ้มค่า

#### ขั้นที่ 35 — Postgres Adapter ที่มี Method เหมือน memory ทุกประการ

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p repository/postgres`

**2. ติดตั้ง Library** `go get github.com/jackc/pgx/v5/pgxpool`

**3. สร้างไฟล์** `repository/postgres/order.go` แล้ววางโค้ดด้านล่างลงไป

**4. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package postgres

import (
	"context"
	"errors"
	"fmt"

	"food-api/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Method มีรูปแบบเหมือน memory.OrderRepo ทุกประการ
// ชั้น Use Case จึงแยกไม่ออกว่ากำลังทำงานกับตัวใด
type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepo {
	return OrderRepo{db: db}
}

func (r OrderRepo) ListMenu() ([]domain.MenuItem, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, name, price FROM menu_items ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var m domain.MenuItem
		if err := rows.Scan(&m.ID, &m.Name, &m.Price); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

func (r OrderRepo) FindMenuItem(id string) (domain.MenuItem, error) {
	var m domain.MenuItem
	err := r.db.QueryRow(context.Background(),
		`SELECT id, name, price FROM menu_items WHERE id = $1`, id,
	).Scan(&m.ID, &m.Name, &m.Price)

	// แปลง Error ของ Driver เป็น Error ของ Domain ที่ชั้นนี้
	// Error ของ pgx จึงไม่รั่วขึ้นไปยังชั้นบน
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.MenuItem{}, domain.ErrNotFound
	}
	return m, err
}

func (r OrderRepo) FindOrder(id string) (domain.Order, error) {
	ctx := context.Background()

	var o domain.Order
	err := r.db.QueryRow(ctx,
		`SELECT id, status FROM orders WHERE id = $1`, id,
	).Scan(&o.ID, &o.Status)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Order{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Order{}, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT m.id, m.name, m.price, l.qty
		FROM order_lines l
		JOIN menu_items m ON m.id = l.menu_id
		WHERE l.order_id = $1
		ORDER BY m.id`, id)
	if err != nil {
		return domain.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var l domain.OrderLine
		if err := rows.Scan(&l.Item.ID, &l.Item.Name, &l.Item.Price, &l.Qty); err != nil {
			return domain.Order{}, err
		}
		o.Lines = append(o.Lines, l)
	}
	return o, rows.Err()
}

// Transaction ของการเขียนข้อมูลอยู่ในชั้นนี้ เพราะเป็นเรื่องของชั้นจัดเก็บข้อมูล
func (r OrderRepo) Save(o domain.Order) error {
	ctx := context.Background()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // จะไม่มีผลหาก Commit สำเร็จไปแล้ว

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (id, status, total) VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status,
		                               total  = EXCLUDED.total`,
		o.ID, string(o.Status), o.Total())
	if err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM order_lines WHERE order_id = $1`, o.ID); err != nil {
		return fmt.Errorf("save order lines: %w", err)
	}

	for _, l := range o.Lines {
		_, err := tx.Exec(ctx,
			`INSERT INTO order_lines (order_id, menu_id, qty) VALUES ($1, $2, $3)`,
			o.ID, l.Item.ID, l.Qty)
		if err != nil {
			return fmt.Errorf("save order lines: %w", err)
		}
	}

	return tx.Commit(ctx)
}
```

- Method มีรูปแบบเหมือน `memory.OrderRepo` ทุกประการ ชั้น Use Case จึงแยกไม่ออกว่ากำลังทำงานกับตัวใด
- `pgx.ErrNoRows` ถูกแปลงเป็น `domain.ErrNotFound` ที่ชั้นนี้ Error ของ Driver จึงไม่รั่วขึ้นไปยังชั้นบน
- `Save` ใช้ Transaction เพราะเขียนสองตารางพร้อมกัน Transaction เป็นเรื่องของชั้นจัดเก็บข้อมูล จึงอยู่ในชั้นนี้ถูกต้องแล้ว
- ยังคงไม่มีกฎทางธุรกิจอยู่ในชั้นนี้ ไม่ปรากฏเลข 100 บาท

<aside>
💡

**สิ่งที่ไม่ถูกแก้แม้แต่ไฟล์เดียวในขั้นนี้** คือ `domain/`, `usecase/`, `delivery/` และ `repository/memory/`

ผลลัพธ์ที่ได้จากการเรียก `curl` เหมือนกับตอนที่ใช้ Memory ทุกตัวอักษร ทั้งที่ชั้นล่างเปลี่ยนจาก Map ไปเป็น Postgres แล้ว นี่คือเหตุผลทั้งหมดของการแบ่งชั้น

</aside>

<aside>
🚨

Connection String อ่านจาก Environment ห้าม Commit ลง Git

หากไม่ได้กำหนด `DATABASE_URL` โปรแกรมจะกลับไปใช้ Memory โดยอัตโนมัติ สังเกตบรรทัดแรกใน Terminal ว่าขึ้น `ใช้ Postgres` หรือ `ใช้ Memory`

</aside>

### 2.5.10 What Survives a Restart

#### ขั้นที่ 36 — สลับที่เก็บข้อมูล โดยไม่แก้ไข domain, usecase และ delivery

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go` ที่สร้างไว้ตั้งแต่ขั้นที่ 11

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

แก้ไขบล็อก `import` ที่ส่วนหัวของไฟล์ ให้เป็นดังนี้

```go
import (
	"context"
	"log"
	"os"

	"food-api/delivery/http"
	"food-api/domain"
	"food-api/repository/memory"
	"food-api/repository/postgres"
	"food-api/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"context"
    	"log"
    	"os"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/repository/postgres"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    	"github.com/jackc/pgx/v5/pgxpool"
    	"github.com/joho/godotenv"
    )
    
    func seedMenu() map[string]domain.MenuItem {
    	return map[string]domain.MenuItem{
    		"m1": {ID: "m1", Name: "Pad Thai", Price: 80},
    		"m2": {ID: "m2", Name: "Tom Yum Goong", Price: 120},
    		"m3": {ID: "m3", Name: "Green Curry", Price: 95},
    		"m4": {ID: "m4", Name: "Thai Iced Tea", Price: 45},
    	}
    }
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึง memory, usecase และ http พร้อมกัน
    // ประกอบจากชั้นในออกมาสู่ชั้นนอก
    func main() {
    	// 0. อ่านค่าตั้งต้นจากไฟล์ .env เข้าสู่ Environment
    	//    การอ่านค่าตั้งต้นเป็นหน้าที่ของชั้นนอกสุดเท่านั้น ชั้นอื่นไม่ควรอ่าน Environment เอง
    	//
    	//    ไม่ตรวจสอบ Error เพราะการไม่มีไฟล์ .env ไม่ใช่ความผิดพลาด
    	//    ผู้ที่ยังไม่ต้องการใช้ฐานข้อมูล ไม่จำเป็นต้องมีไฟล์นี้เลย
    	//
    	//    ค่าที่กำหนดมาจาก Terminal จะชนะค่าในไฟล์ .env เสมอ
    	//    จึงยังสั่ง DATABASE_URL= go run . เพื่อบังคับกลับไปใช้ Memory ได้
    	_ = godotenv.Load()
    
    	// 1. ที่เก็บข้อมูล
    	//    เลือกจาก Environment ไม่ใช่จากการแก้โค้ด
    	//    เมื่อกำหนด DATABASE_URL ระบบจะเปลี่ยนไปใช้ Postgres ทันที
    	var repo usecase.OrderRepository
    
    	if url := os.Getenv("DATABASE_URL"); url != "" {
    		pool, err := pgxpool.New(context.Background(), url)
    		if err != nil {
    			log.Fatal(err)
    		}
    		defer pool.Close()
    
    		if err := pool.Ping(context.Background()); err != nil {
    			log.Fatal("เชื่อมต่อฐานข้อมูลไม่สำเร็จ: ", err)
    		}
    
    		log.Println("ใช้ Postgres")
    		repo = postgres.NewOrderRepo(pool) // <-- สลับที่เก็บข้อมูลตรงนี้
    	} else {
    		log.Println("ใช้ Memory")
    		repo = memory.NewOrderRepo(seedMenu())
    	}
    
    	// 2. กฎ
    	//    domain, usecase และ delivery ไม่ถูกแก้ไขแม้แต่ไฟล์เดียว
    	uc := usecase.NewOrderUseCase(repo)
    
    	// 3. ชั้นขนส่ง
    	h := http.NewHandler(uc)
    
    	app := fiber.New()
    
    	// Middleware ครอบทุก Request โดยเขียนไว้ที่เดียว (Decorator Pattern)
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	api := app.Group("/api")
    	api.Get("/menu", h.ListMenu)
    	api.Post("/orders", h.PlaceOrder)
    	api.Get("/orders/:id", h.GetOrder)
    	api.Get("/orders/:id/total", h.OrderTotal)
    	api.Patch("/orders/:id/status", h.SetStatus)
    	api.Post("/orders/:id/cancel", h.CancelOrder)
    
    	log.Fatal(app.Listen(":3000"))
    }
    ```
    

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -i localhost:3000/api/menu
```

```bash
HTTP/1.1 200 OK
Content-Type: application/json

[{"id":"m1","name":"Pad Thai","price":80},{"id":"m2","name":"Tom Yum Goong","price":120},{"id":"m3","name":"Green Curry","price":95},{"id":"m4","name":"Thai Iced Tea","price":45}]
```

คำสั่งที่ 2

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"m1","qty":2}'
```

```bash
HTTP/1.1 201 Created
Content-Type: application/json

{"id":"o1","status":"pending","total":160,"lines":[{"menu_id":"m1","name":"Pad Thai","price":80,"qty":2}]}
```

คำสั่งที่ 3

```bash
curl -i localhost:3000/api/orders/o1
```

```bash
HTTP/1.1 200 OK
Content-Type: application/json

{"id":"o1","status":"pending","total":160,"lines":[{"menu_id":"m1","name":"Pad Thai","price":80,"qty":2}]}
```

<aside>
💡

Named Volume มีอายุยืนกว่า Container จะปิดเปิดกี่ครั้งตลอดสัปดาห์ ข้อมูลก็ยังอยู่ครบ มีเพียง `docker compose down -v` เท่านั้นที่ลบทิ้ง

สังเกตว่าหลัง `down -v` ฐานข้อมูลกลับไปว่างเปล่าจนไม่มีตารางใดเลย ต้องรัน `./cmd/migrate` ใหม่ Schema จึงจะกลับมา

และเมื่อรันแล้ว `menu_items` กลับมามี 4 แถวเพราะอยู่ใน `001_init.sql` ส่วน `orders` ว่างเปล่า เพราะเป็นข้อมูลที่ผู้ใช้สร้างขึ้นภายหลัง

</aside>

<aside>
🚨

`docker compose down -v` เป็นทั้งปุ่มรีเซ็ตและเป็นอุบัติเหตุที่พบบ่อยที่สุด ในวันแรกที่ยังไม่มีข้อมูลสำคัญ ใช้ได้โดยอิสระ แต่หลังจากนั้นต้องพิจารณาให้รอบคอบก่อนใช้ทุกครั้ง

ในงานจริง ขั้นตอน Deploy จะเรียก `./cmd/migrate` ทุกครั้งก่อนเริ่มโปรแกรม เพื่อให้ Schema เป็นรุ่นล่าสุดเสมอ

</aside>

#### ขั้นที่ 37 — อะไรรอดจากการปิดเปิด และอะไรหายถาวร

<aside>
📄

**1. สร้างออเดอร์ไว้หนึ่งรายการก่อน** เพื่อให้มีข้อมูลไว้สังเกต

**2. ลอง** `docker compose down` แล้ว `docker compose up -d` จากนั้นตรวจสอบข้อมูลอีกครั้ง

**3. ลอง** `docker compose down -v` แล้ว `docker compose up -d` จากนั้นตรวจสอบข้อมูลอีกครั้ง

**4. รัน Migration ใหม่** แล้วสังเกตว่าอะไรกลับมา อะไรหายถาวร

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ผลลัพธ์**

คำสั่งที่ 1

```bash
docker compose exec db psql -U app -d food -c "SELECT id,status,total FROM orders;"
```

```bash
 id | status  | total
----+---------+--------
 o1 | pending | 160.00
(1 row)
```

คำสั่งที่ 2

```bash
docker compose down          # ลบ container แต่ volume ยังอยู่
```

```bash
Container food-api-db-1   Removed
Network food-api_default  Removed
```

คำสั่งที่ 3

```bash
docker compose up -d
docker compose exec db psql -U app -d food -c "SELECT id,status,total FROM orders;"
```

```bash
 id | status  | total
----+---------+--------
 o1 | pending | 160.00
(1 row)

   ← ข้อมูลยังอยู่ครบ
```

คำสั่งที่ 4

```bash
docker compose down -v       # ลบ volume ด้วย ข้อมูลหายถาวร
```

```bash
Container food-api-db-1   Removed
Volume food-api_pgdata    Removed
Network food-api_default  Removed
```

คำสั่งที่ 5

```bash
docker compose up -d
docker compose exec db psql -U app -d food -c "\dt"
```

```bash
Did not find any relations.

   ← ไม่เหลือแม้แต่ตารางเดียว
```

คำสั่งที่ 6

```bash
go run ./cmd/migrate
```

```bash
ใช้งาน 001_init แล้ว
ใช้งาน 002_add_note_to_orders แล้ว
ใช้งาน 003_index_orders_status แล้ว
ใช้งาน migration ใหม่ทั้งหมด 3 ไฟล์
```

คำสั่งที่ 7

```bash
docker compose exec db psql -U app -d food -c "SELECT count(*) FROM menu_items;" -c "SELECT count(*) FROM orders;"
```

```bash
 count
-------
     4
(1 row)

 count
-------
     0
(1 row)

   ← เมนูกลับมาจาก 001 แต่ออเดอร์ของผู้ใช้หายถาวร
```

### 2.5.11 Troubleshooting

| อาการ | สาเหตุและทางแก้ |
| --- | --- |
| Port 5432 ถูกใช้อยู่ | มี Postgres ตัวอื่นรันอยู่แล้ว ให้ใช้ตัวนั้นแทนโดยแก้ค่า `POSTGRES_*` ใน `.env` ให้ชี้ไปที่ตัวนั้น หรือหากยังต้องการใช้ Docker ให้แก้ `POSTGRES_PORT` ใน `.env` เป็น `5433` แก้บรรทัดเดียวพอ เพราะ `DATABASE_URL` ประกอบขึ้นจากค่านั้นอยู่แล้ว |
| ต่อ Docker ไม่ได้ | Docker Desktop ยังไม่ได้เปิด หรือยังไม่ได้เปิด WSL Integration |
| `connection refused` | โปรแกรมเริ่มก่อนที่ Postgres จะบูตเสร็จ ให้รอจน `docker compose ps` ขึ้น `(healthy)` ก่อน |
| แก้รหัสผ่านใน `.env` แล้วแต่ยังต่อไม่ได้ | Volume เดิมยังจำรหัสผ่านชุดแรกไว้ เพราะ Postgres ตั้งรหัสผ่านเฉพาะตอนสร้าง Volume ครั้งแรกเท่านั้น ต้อง `docker compose down -v` แล้ว `up -d` ใหม่ ซึ่งจะลบข้อมูลทั้งหมด |
| `relation "orders" does not exist` | ยังไม่ได้รัน Migration ให้รัน `go run ./cmd/migrate` ก่อนเริ่มโปรแกรม |
| เพิ่มไฟล์ `.sql` ใหม่แล้วแต่ไม่มีอะไรเปลี่ยน | ลืมรัน `./cmd/migrate` หลังเพิ่มไฟล์ หรือรันผิดโฟลเดอร์ ต้องรันจากโฟลเดอร์ `food-api` |
| Migration รันแล้วแต่ Schema ไม่ตรงกับเพื่อน | มีคนแก้ไฟล์ที่รันไปแล้ว ซึ่งจะไม่ถูกรันซ้ำ ให้เพิ่มไฟล์ใหม่แทนการแก้ไฟล์เดิมเสมอ |
| ข้อมูลหายไป | มีผู้รัน `docker compose down -v` ซึ่งเป็นคำสั่งเดียวที่ลบข้อมูล |
| `address already in use` ตอน `go run .` | มีโปรแกรมอื่นใช้ Port `3000` อยู่ ให้ปิดโปรแกรมนั้น หรือแก้เลขใน `app.Listen` |
| `ไม่พบ DATABASE_URL` | ยังไม่ได้สร้างไฟล์ `.env` ให้สั่ง `cp .env.example .env` หรือรันคำสั่งผิดโฟลเดอร์ เพราะ `.env` ถูกอ่านจากโฟลเดอร์ที่รันคำสั่ง ต้องรันในโฟลเดอร์ `food-api` |
| `database "food" does not exist` `(SQLSTATE 3D000)` | ต่อเข้าเซิร์ฟเวอร์ได้แล้ว แต่ยังไม่มีฐานข้อมูลชื่อนั้น พบในผู้ที่ใช้ Postgres ของตนเอง เพราะไฟล์ Migration สร้างเฉพาะตาราง ไม่ได้สร้างฐานข้อมูล ให้สั่ง `createdb food` ก่อน |
| `password authentication failed` `(SQLSTATE 28P01)` | ค่า `POSTGRES_USER` หรือ `POSTGRES_PASSWORD` ใน `.env` ไม่ตรงกับฐานข้อมูลที่กำลังทำงานอยู่ ให้แก้ที่สองบรรทัดนั้น แล้ว `DATABASE_URL` จะถูกต้องตามไปเอง |
| `dial error` ... `connection refused` | ไม่มีสิ่งใดฟังอยู่ที่ Port นั้น ให้ตรวจ `POSTGRES_PORT` ใน `.env` ว่าตรงกับ Port ที่ Postgres เปิดอยู่จริงหรือไม่ เพราะทั้งเลขหน้าเครื่องหมาย `:` ในไฟล์ Compose และเลข Port ใน `DATABASE_URL` ดึงมาจากค่านี้ค่าเดียว |
| แก้ `.env` แล้วแต่ไม่มีผล | มีการกำหนดตัวแปรเดียวกันไว้ใน Terminal อยู่แล้ว ซึ่งชนะค่าในไฟล์ `.env` เสมอ ให้เปิด Terminal ใหม่ หรือรันคำสั่งผิดโฟลเดอร์ |

คำสั่งที่ใช้ไล่หาสาเหตุ

```bash
# 1. ใครใช้ port 5432 อยู่
lsof -i :5432                    # macOS และ Linux
netstat -ano | findstr :5432     # Windows

brew services stop postgresql@16    # macOS
sudo systemctl stop postgresql      # Linux

# 2. Docker เปิดอยู่หรือไม่
docker info

# 3. ฐานข้อมูลขึ้นแล้วจริงหรือยัง
docker compose ps
docker compose logs db | tail -30
#   มองหาบรรทัด: database system is ready to accept connections

# 4. พิสูจน์ว่าคุยกับฐานข้อมูลได้
docker compose exec db psql -U app -d food -c "\dt"
#   หากคำสั่งนี้ทำงานได้ แสดงว่าปัญหาอยู่ที่โค้ด Go ไม่ใช่ที่ Docker
```

### 2.5.12 Designing a Schema That Can Change

วิชาฐานข้อมูลสอนให้เก็บ Requirement ให้ครบ วาด ER Diagram ให้เสร็จ แล้ว Normalize ทุกตารางก่อนเขียนโค้ด วิธีนี้ถูกต้องเมื่อความต้องการนิ่งแล้ว แต่ในงานจริงความต้องการเปลี่ยนเกือบทุก Release

| ในวิชาฐานข้อมูล | ในงานจริง |
| --- | --- |
| เก็บ Requirement ให้ครบก่อน | เริ่มจาก Use Case เดียว และตารางที่จำเป็นต้องใช้ |
| วาด ER Diagram ทั้งระบบก่อนเขียนโค้ด | Schema โตไปพร้อมกับโค้ด |
| Normalize ทุกอย่างจนถึง 3NF | Normalize เฉพาะจุดที่คุ้ม |
| ออกแบบเสร็จแล้วแช่แข็ง | ทุก Release มีสิทธิ์แก้ Schema |
| วัดผลที่ความครบถ้วน | วัดผลที่ความสามารถในการเปลี่ยน |

เมื่อยังไม่ทราบว่าข้อมูลจะมีรูปแบบอย่างไร ให้แยกเป็นแกนที่นิ่งกับขอบที่ยังเปลี่ยน

```sql
CREATE TABLE orders (
    -- แกนที่นิ่งแล้ว: ใช้ค้นหา ใส่ Constraint และทำ Index
    id       TEXT PRIMARY KEY,
    status   TEXT NOT NULL,
    total    NUMERIC(10,2) NOT NULL,
    created  TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- ขอบที่ยังไม่นิ่ง: รูปแบบยังไม่แน่นอน
    attributes JSONB NOT NULL DEFAULT '{}'
);

-- เลื่อนขั้นเป็นคอลัมน์จริงเมื่อรูปแบบนิ่งแล้ว
ALTER TABLE orders
    ADD COLUMN delivery_note TEXT;
```

- สิ่งที่ต้องใช้ Filter, Join หรือใส่ Constraint ต้องเป็นคอลัมน์จริง JSONB ไม่ใช่ที่สำหรับซ่อนข้อมูล
- JSONB เหมาะกับค่าที่ต่างกันไปในแต่ละร้าน Flag ของการทดลอง หรือฟีลด์ที่ยังตั้งชื่อไม่ได้
- เริ่มจาก JSONB แล้วย้ายเป็นคอลัมน์เมื่อรูปแบบนิ่งแล้ว

หลักที่ทำให้ Schema เปลี่ยนต่อได้เรื่อย ๆ

1. **เพิ่มได้ แต่อย่าลบ** การเพิ่มคอลัมน์แบบ Nullable ปลอดภัยกับตารางที่มีข้อมูลอยู่แล้ว ส่วนการลบคอลัมน์ไม่ปลอดภัย
2. **Expand แล้วค่อย Contract** เพิ่มคอลัมน์ใหม่ เขียนลงทั้งสองที่ ย้ายข้อมูลเก่า สลับฝั่งอ่าน แล้วจึงลบของเดิม รวมสี่รอบ Deploy โดยระบบไม่ต้องหยุด
3. **ห้ามเปลี่ยนชื่อในรอบเดียว** การ Rename ทำให้โค้ดที่ทำงานอยู่ล้มเหลวทันทีที่คำสั่งถูกเรียกใช้ ให้เพิ่มใหม่ ย้าย แล้วค่อยลบ
4. **อย่ารีบ Normalize** อย่าเพิ่งแยกตารางที่ยังไม่เคยสร้างปัญหา Join ที่ไม่จำเป็นต้องจ่ายทุกวัน
5. **Key ต้องไม่มีความหมาย** `id` ที่เป็น Surrogate Key ไม่มีวันเปลี่ยน แต่อีเมลหรือเบอร์โทรเปลี่ยนแน่นอน
6. **Domain คือความจริง** Struct ใน Go เป็นตัวกำหนดแบบจำลอง ส่วน Schema เป็นเพียงที่เก็บ นี่คือเหตุผลที่ชั้น `domain` ไม่มี SQL อยู่เลย

## 2.6 Example Problems

โจทย์ทั้งสามข้อทำต่อจากโค้ดที่สร้างเสร็จแล้วในหัวข้อ 2.4 และ 2.5 ได้ทันที

### 2.6.1 Cancel Endpoint ที่เขียนผิดหลัก

จงพิจารณาโค้ดต่อไปนี้ แล้วตอบคำถามโดยไม่ต้องรัน

```go
func (h Handler) CancelOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	row := h.db.QueryRow("SELECT status FROM orders WHERE id = $1", id)
	var status string
	row.Scan(&status)

	if status == "cooking" || status == "delivering" {
		return c.Status(409).JSON(msg("cannot cancel"))
	}

	h.db.Exec("UPDATE orders SET status = 'cancelled' WHERE id = $1", id)
	return c.SendStatus(204)
}
```

1. โค้ดนี้ทำงานถูกต้องหรือไม่ และผิดหลัก Dependency Rule ข้อใด
2. หากในอนาคตมีสถานะใหม่ชื่อ `packing` ที่ยกเลิกไม่ได้ ต้องแก้ไขกี่จุด และมีความเสี่ยงที่จะหลงลืมที่จุดใดบ้าง
3. จงเปรียบเทียบกับ `CancelOrder` ที่เขียนไว้ที่ **ขั้นที่ 22** แล้วอธิบายว่ากฎการยกเลิกย้ายไปอยู่ที่ใด และเพราะเหตุใดจึงดีกว่า

<aside>
⚠️

**Challenge: จงเพิ่มสถานะ `packing` ระหว่าง `confirmed` กับ `cooking` เข้าไปในระบบจริง โดยแก้ไฟล์ให้น้อยที่สุด แล้วนับว่าแก้ไปกี่ไฟล์ กี่บรรทัด**

</aside>

### 2.6.2 Endpoint ใหม่ — ยกเลิกทั้งวัน

จงเพิ่ม Endpoint `GET /api/orders?status=pending` ที่คืนรายการออเดอร์ทั้งหมดที่มีสถานะตามที่ระบุ ตามเงื่อนไขต่อไปนี้

1. เพิ่ม Method `ListOrders(status domain.Status) ([]domain.Order, error)` ให้ `OrderUseCase`
2. ขยาย `OrderRepository` และทำให้ทั้ง `memory` และ `postgres` ครบตาม Interface ใหม่
3. เขียน Handler ที่อ่าน Query Parameter เรียก Use Case แปลง Error แล้วคืน JSON
4. ห้ามมี SQL และห้ามมี `fiber.Ctx` อยู่ในชั้น Use Case และห้ามมีการกรองข้อมูลอยู่ในชั้น Handler

```go
// GET /api/orders?status=pending
// [{"id":"o1","status":"pending","total":160,"lines":[...]}]

func (h Handler) ListOrders(c *fiber.Ctx) error {
	// Write your code here
}
```

<aside>
⚠️

**Challenge: หากส่ง `status` ที่ไม่มีอยู่จริง เช่น `?status=banana` ควรคืน Status Code ใด และควรตรวจสอบที่ชั้นใด จงเพิ่ม Sentinel Error ใหม่แล้วแก้เฉพาะ `toHTTP` เท่านั้น**

</aside>

### 2.6.3 Schema That Can Change

ระบบเปิดใช้งานจริงแล้ว มีออเดอร์ของลูกค้าอยู่ในฐานข้อมูล และ Migration ล่าสุดคือ `003`

ร้านค้าแจ้งความต้องการใหม่เข้ามาสามข้อ

1. ทุกออเดอร์ต้องเก็บ **ชื่อผู้รับ** และ **เบอร์โทร** เพราะไรเดอร์ต้องใช้ติดต่อ และแอดมินต้องค้นหาออเดอร์จากเบอร์โทรได้
2. บางร้านที่มีหน้าร้านอยากเก็บ **หมายเลขโต๊ะ** แต่ร้านที่ขายเฉพาะเดลิเวอรีไม่มีโต๊ะ
3. ทีมการตลาดอยากทดลองเก็บ **Coupon Code** แต่ยังไม่แน่ใจว่าจะใช้จริงหรือไม่ และยังไม่ต้องนำมาค้นหา

จงตอบว่าแต่ละข้อควรเป็น **คอลัมน์จริง** หรือควรอยู่ใน `attributes` ที่เป็น JSONB พร้อมเหตุผล

แล้วเขียนไฟล์ `migrations/004_....sql` สำหรับข้อที่เลือกให้เป็นคอลัมน์ และรันด้วย `./cmd/migrate`

<aside>
🚨

ห้ามแก้ไฟล์ `001` ถึง `003` เพราะรันไปแล้ว การแก้ไฟล์ที่รันแล้วจะไม่มีผลกับฐานข้อมูลที่ Migration ผ่านไปแล้ว

และจะทำให้ Schema ของแต่ละเครื่องไม่ตรงกัน

</aside>

<aside>
⚠️

**Challenge: หากภายหลัง Coupon Code ถูกใช้จริงและต้องนำมา Filter จงอธิบายขั้นตอน Expand แล้ว Contract ทั้งสี่รอบ Deploy ที่ทำให้ย้ายจาก JSONB ไปเป็นคอลัมน์ได้โดยระบบไม่ต้องหยุด**

</aside>

## 2.7 Solutions

<aside>
🚨

ให้ลองทำด้วยตนเองจนสุดทางก่อน แล้วจึงเปิดดูเฉลย การอ่านเฉลยที่ถูกต้องไม่ได้ทำให้เขียนโค้ดเองได้

</aside>

### 2.7.1 Cancel Endpoint ที่เขียนผิดหลัก

- กดเพื่อดูเฉลย
    
    **ข้อ 1** โค้ดทำงานได้ถูกต้องในแง่ผลลัพธ์ แต่ผิดหลัก Dependency Rule เพราะชั้น Delivery เรียก SQL โดยตรง เท่ากับชั้นนอกสุดกระโดดข้ามชั้น Use Case และชั้น Domain ลงไปคุยกับชั้นนอกอีกฝั่งหนึ่ง
    
    ผลที่ตามมาคือกฎการยกเลิกซึ่งเป็นกฎทางธุรกิจ ไปอยู่ในไฟล์ที่ควรมีหน้าที่แปลง HTTP เท่านั้น และยังผูกติดกับ Postgres จนทดสอบไม่ได้หากไม่มีฐานข้อมูล
    
    นอกจากนี้ยังละเลย Error จาก `row.Scan` หากไม่มีออเดอร์นั้นอยู่จริง `status` จะเป็นสตริงว่าง ซึ่งไม่ตรงกับ `cooking` หรือ `delivering` โค้ดจึงเดินต่อไปสั่ง `UPDATE` แถวที่ไม่มีอยู่ แล้วคืน `204` ราวกับว่าสำเร็จ
    
    **ข้อ 2** ต้องแก้ทุกที่ที่มีการเทียบสถานะแบบนี้ ซึ่งกระจายอยู่ตาม Handler ทุกตัวที่เกี่ยวข้อง ในระบบจริงมักลืมที่ Endpoint ที่ไม่ได้ใช้บ่อย เช่นหน้าจัดการของแอดมิน หรือ Job ที่ทำงานตามเวลา ปัญหาคือไม่มีจุดใดที่ระบุได้ว่า *กฎการยกเลิกอยู่ตรงนี้*
    
    **ข้อ 3** กฎย้ายไปอยู่ใน `domain` เพราะเป็นข้อเท็จจริงเกี่ยวกับ `Order` ส่วนการตัดสินใจว่าจะทำอย่างไรต่ออยู่ใน `usecase` และ Handler เหลือเพียงการแปลง
    
    ```go
    // domain/order.go — กฎอยู่ติดกับ Type
    func (o Order) CanCancel() bool {
    	return o.Status == Pending || o.Status == Confirmed
    }
    
    // usecase/order.go — ตัดสินใจและบันทึก
    func (u OrderUseCase) Cancel(id string) error {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return fmt.Errorf("cancel order %s: %w", id, err)
    	}
    
    	if !order.CanCancel() {
    		return domain.ErrNotCancelable
    	}
    
    	order.Status = domain.Cancelled
    	return u.repo.Save(order)
    }
    
    // delivery/http/order.go — แปลงอย่างเดียว
    func (h Handler) CancelOrder(c *fiber.Ctx) error {
    	if err := h.uc.Cancel(c.Params("id")); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    ```
    
    **Challenge** — เพิ่มสถานะ `packing`
    
    แก้ไฟล์เดียว คือ `domain/order.go` และแก้บรรทัดเดียว
    
    ```go
    var flow = []Status{Pending, Confirmed, Packing, Cooking, Delivering, Completed}
    ```
    
    พร้อมเพิ่มค่าคงที่อีกหนึ่งบรรทัดในบล็อก `const`
    
    ```go
    Packing Status = "packing"
    ```
    
    `CanCancel` ไม่ต้องแก้เลย เพราะเขียนไว้ว่า *ยกเลิกได้เฉพาะ `pending` และ `confirmed`* ไม่ได้เขียนว่า *ยกเลิกไม่ได้หากเป็น cooking หรือ delivering*
    
    ความต่างของการเขียนกฎแบบ Allow List กับ Deny List ปรากฏผลชัดเจนตรงนี้ หากเขียนแบบ Deny List อย่างในโจทย์ การเพิ่มสถานะใหม่จะทำให้ `packing` ยกเลิกได้ทันทีโดยไม่มีผู้ใดสังเกตเห็น
    
    `delivery` และ `repository` ไม่ถูกแก้ไขเลยแม้แต่บรรทัดเดียว
    

### 2.7.2 Endpoint ใหม่ — ยกเลิกทั้งวัน

- กดเพื่อดูเฉลย
    
    ขยาย Interface และเพิ่ม Method ในชั้น Use Case
    
    ```go
    // usecase/order.go
    type OrderRepository interface {
    	ListMenu() ([]domain.MenuItem, error)
    	FindMenuItem(id string) (domain.MenuItem, error)
    	FindOrder(id string) (domain.Order, error)
    	ListOrders(status domain.Status) ([]domain.Order, error)
    	Save(o domain.Order) error
    }
    
    func (u OrderUseCase) ListOrders(status domain.Status) ([]domain.Order, error) {
    	return u.repo.ListOrders(status)
    }
    ```
    
    ฝั่ง Repository การกรองอยู่ที่นี่ เพราะเป็นเรื่องของการดึงข้อมูล
    
    ```go
    // repository/memory/order.go
    func (r OrderRepo) ListOrders(status domain.Status) ([]domain.Order, error) {
    	r.mu.RLock()
    	defer r.mu.RUnlock()
    
    	out := make([]domain.Order, 0)
    	for _, o := range r.orders {
    		if o.Status == status {
    			out = append(out, o)
    		}
    	}
    	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
    	return out, nil
    }
    ```
    
    Handler ทำเพียงแปลงเข้าและแปลงออก
    
    ```go
    // delivery/http/order.go
    func (h Handler) ListOrders(c *fiber.Ctx) error {
    	orders, err := h.uc.ListOrders(domain.Status(c.Query("status", "pending")))
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	out := make([]orderResponse, 0, len(orders))
    	for _, o := range orders {
    		out = append(out, toResponse(o))
    	}
    	return c.JSON(out)
    }
    ```
    
    แล้วเพิ่ม Route หนึ่งบรรทัดใน `main.go`
    
    ```go
    api.Get("/orders", h.ListOrders)
    ```
    
    <aside>
    🚨
    
    การขยาย Interface ครั้งนี้ทำให้มี **4 ที่** ที่ต้องเพิ่ม `ListOrders` ไม่ใช่ 2 ที่
    
    `memory.OrderRepo` และ `postgres.OrderRepo` ตามที่คาดไว้
    
    และ `stubRepo` กับ `failingRepo` ใน `usecase/fakes_test.go` ด้วย มิฉะนั้น `go test` จะไม่ผ่าน
    
    `*stubRepo does not implement usecase.OrderRepository (missing method ListOrders)`
    
    นี่คือราคาที่ต้องจ่ายเวลาขยาย Interface และเป็นเหตุผลที่ควรทำ Interface ให้เล็กที่สุดเท่าที่จำเป็น
    
    </aside>
    
    ผลลัพธ์ที่ได้
    
    ```bash
    curl "localhost:3000/api/orders?status=pending"
    ```
    
    ```bash
    [{"id":"o1","status":"pending","total":160,"lines":[{"menu_id":"m1","name":"Pad Thai","price":80,"qty":2}]}]
    ```
    
    ```bash
    curl "localhost:3000/api/orders?status=confirmed"
    ```
    
    ```bash
    [{"id":"o2","status":"confirmed","total":120,"lines":[{"menu_id":"m2","name":"Tom Yum Goong","price":120,"qty":1}]}]
    ```
    
    **Challenge** — `?status=banana`
    
    ควรคืน `409` หรือ `400` ก็ได้ แต่ต้องตรวจที่ชั้น **Use Case** ไม่ใช่ที่ Handler เพราะ *สถานะใดมีอยู่จริงบ้าง* เป็นความรู้ของโดเมน ไม่ใช่ความรู้ของ HTTP
    
    ```go
    // domain/order.go
    func (s Status) Valid() bool {
    	switch s {
    	case Pending, Confirmed, Cooking, Delivering, Completed, Cancelled:
    		return true
    	}
    	return false
    }
    ```
    
    แล้วให้ `ListOrders` ในชั้น Use Case เรียกใช้ก่อนส่งต่อไปที่ Repository
    
    ```go
    // usecase/order.go
    func (u OrderUseCase) ListOrders(status domain.Status) ([]domain.Order, error) {
    	if !status.Valid() {
    		return nil, domain.ErrInvalidStatus
    	}
    	return u.repo.ListOrders(status)
    }
    ```
    
    `ErrInvalidStatus` ถูกจัดการอยู่แล้วใน `toHTTP` จึงได้ `409` มาทันทีโดยไม่ต้องเขียนเพิ่มในชั้น Delivery เลย นี่คือผลของการรวมการตัดสินใจไว้ที่เดียว
    
    ```bash
    curl -i "localhost:3000/api/orders?status=banana"
    ```
    
    ```bash
    HTTP/1.1 409 Conflict
    {"error":"invalid status transition"}
    ```
    
    <aside>
    ⚠️
    
    สังเกตข้อความที่ได้กลับมาคือ `invalid status transition` ซึ่ง **ไม่ตรงกับความจริง** เพราะผู้ใช้ไม่ได้พยายามเปลี่ยนสถานะ เพียงส่งชื่อสถานะที่ไม่มีอยู่จริงเข้ามา
    
    การหยิบ Sentinel Error ตัวเดิมมาใช้ซ้ำเพราะ *คืน 409 ให้ทันที* เป็นทางลัดที่ทำให้ข้อความผิดความหมาย ทางที่ถูกกว่าคือประกาศ `ErrUnknownStatus` ขึ้นมาใหม่ แล้วเพิ่มอีกหนึ่ง `case` ใน `toHTTP` ให้คืน `400` พร้อมข้อความที่ตรงกับสิ่งที่เกิดขึ้นจริง
    
    คำถามที่ควรถามทุกครั้งที่จะใช้ Sentinel Error ตัวเดิมซ้ำคือ *ข้อความนี้อธิบายสิ่งที่เกิดขึ้นจริงหรือไม่* ไม่ใช่ *Error ตัวนี้คืน Status Code ที่ต้องการหรือไม่*
    
    </aside>
    

### 2.7.3 Schema That Can Change

- กดเพื่อดูเฉลย
    
    
    | ความต้องการ | ควรเก็บที่ไหน | เหตุผล |
    | --- | --- | --- |
    | ชื่อผู้รับ และเบอร์โทร | คอลัมน์จริง | ทุกออเดอร์ต้องมี ความหมายชัดเจน และโจทย์ระบุว่าต้องใช้ค้นหา สิ่งที่ต้องใช้ค้นหาต้องเป็นคอลัมน์เพื่อทำ Index ได้ |
    | หมายเลขโต๊ะ | `attributes` | มีเฉพาะบางร้าน หากทำเป็นคอลัมน์จะว่างเป็นส่วนใหญ่ และยังไม่จำเป็นต้องใช้ค้นหา |
    | Coupon Code | `attributes` | ยังเป็นการทดลอง หากเลิกใช้ก็ลบทิ้งได้โดยไม่ต้องแก้ไข Schema |
    
    ไฟล์ Migration ใหม่ ใช้เลข `004` ต่อจากของเดิม
    
    ```sql
    -- migrations/004_add_receiver_to_orders.sql
    
    -- ต้องเป็น Nullable เพราะออเดอร์เก่าที่มีอยู่แล้วไม่มีค่านี้
    ALTER TABLE orders ADD COLUMN receiver_name  TEXT;
    ALTER TABLE orders ADD COLUMN receiver_phone TEXT;
    
    -- โจทย์ระบุว่าแอดมินต้องค้นหาจากเบอร์โทร จึงต้องมี Index
    CREATE INDEX idx_orders_receiver_phone ON orders(receiver_phone);
    ```
    
    แล้วรัน
    
    ```bash
    go run ./cmd/migrate
    ```
    
    ```bash
    ใช้งาน 004_add_receiver_to_orders แล้ว
    ใช้งาน migration ใหม่ทั้งหมด 1 ไฟล์
    ```
    
    <aside>
    💡
    
    หมายเลขโต๊ะและ Coupon Code ไม่ต้องเขียน Migration เลย เพราะคอลัมน์ `attributes` มีอยู่แล้วตั้งแต่ `002`
    
    เขียนลงไปได้ทันทีในรูปแบบ `{"table_no": 12, "coupon_code": "NEWYEAR"}` โดยไม่ต้องแตะ Schema
    
    นี่คือประโยชน์ของการเผื่อช่อง JSONB ไว้ตั้งแต่ต้น
    
    </aside>
    
    <aside>
    🚨
    
    หากใส่ `NOT NULL` โดยไม่มี `DEFAULT` คำสั่งจะล้มเหลวทันทีบนตารางที่มีข้อมูล
    
    `ERROR: column "receiver_name" of relation "orders" contains null values`
    
    หากจำเป็นต้องบังคับ `NOT NULL` จริง ๆ ต้องทำสามรอบ คือเพิ่มแบบ Nullable ก่อน แล้ว Backfill ค่าเดิม จากนั้นจึงค่อย `SET NOT NULL`
    
    </aside>
    
    **Challenge** — ย้าย Coupon Code จาก JSONB ไปเป็นคอลัมน์ โดยระบบไม่หยุด
    
    1. **Expand** เพิ่มคอลัมน์ `coupon_code TEXT` แบบ Nullable ใน `005` โดยยังไม่มีโค้ดใดอ่านจากคอลัมน์นี้
    2. **Write both** แก้โค้ดให้เขียนทั้งลง `attributes` และลงคอลัมน์ใหม่พร้อมกัน Deploy แล้วปล่อยให้ข้อมูลใหม่ไหลเข้าทั้งสองที่
    3. **Backfill แล้วสลับฝั่งอ่าน** รัน `UPDATE` ย้ายข้อมูลเก่าจาก JSONB ลงคอลัมน์ ตรวจว่าครบ แล้วจึงแก้โค้ดให้อ่านจากคอลัมน์ Deploy อีกครั้ง
    4. **Contract** เมื่อไม่มีโค้ดใดอ่านจาก `attributes` แล้ว จึงหยุดเขียนลง JSONB และลบ Key นั้นทิ้ง
    
    ```sql
    -- migrations/006_backfill_coupon_code.sql  (ขั้นตอนที่ 3)
    UPDATE orders
    SET coupon_code = attributes ->> 'coupon_code'
    WHERE coupon_code IS NULL
      AND attributes ? 'coupon_code';
    ```
    
    <aside>
    💡
    
    เหตุผลที่ต้องแยกเป็นสี่รอบ เพราะระหว่าง Deploy จะมีทั้งโค้ดเวอร์ชันเก่าและใหม่ทำงานพร้อมกันอยู่ช่วงหนึ่ง ทุกขั้นตอนจึงต้องทำงานได้ถูกต้องกับโค้ดทั้งสองเวอร์ชัน
    
    สังเกตว่าทั้งสี่รอบเป็นการ **เพิ่ม** ทั้งหมด ไม่มีรอบใดที่ลบของที่ยังมีคนใช้อยู่
    
    </aside>
    

## 2.8 สรุปสิ่งที่ได้ทำวันนี้

| ขั้นที่ | ได้สิ่งใด |
| --- | --- |
| 1 – 4 | เขียน Fiber ได้ รับ Request อ่าน JSON เลือก Status Code และรวมการแปลง Error ไว้ที่เดียว |
| 5 – 11 | วางโครง 4 ชั้นให้ครบ แล้วทำให้ `GET /api/menu` วิ่งผ่านทุกชั้นได้จริง |
| 12 – 19 | ใส่ `Order` กฎยอดขั้นต่ำ และการแปลง Error เป็น `404`, `409`, `422` |
| 20 – 23 | กฎการยกเลิกและการเดินสถานะ ซึ่งเป็นเหตุผลทั้งหมดของการมีชั้น Use Case |
| 24 – 27 | เพิ่ม Endpoint ใหม่ทั้งเส้น และเขียน Test ที่ไม่ต้องใช้ฐานข้อมูล |
| 28 | แยกค่าตั้งต้นออกจากโค้ดไปไว้ใน `.env` และกัน `.env` ไม่ให้ขึ้น Git |
| 29 – 34 | เปิด Postgres ด้วย Docker Compose และเขียน Migration 3 ไฟล์ที่รันบนฐานข้อมูลที่มีข้อมูลอยู่จริง |
| 35 – 37 | ย้ายข้อมูลไป Postgres โดยไม่แก้ `domain`, `usecase` และ `delivery` แม้แต่บรรทัดเดียว |

<aside>
💡

ข้อสรุปของทั้งวัน อยู่ที่ **ขั้นที่ 36** หากย้ายจาก Map ไปเป็น Postgres ได้โดยแก้ไฟล์เดียว แสดงว่าเส้นแบ่งชั้นถูกกำหนดไว้อย่างเหมาะสมแล้ว

หากวันใดวันหนึ่ง การเปลี่ยนฐานข้อมูลทำให้ต้องเปิดไฟล์ในโฟลเดอร์ `domain` แสดงว่ามีบางอย่างรั่วข้ามชั้นไปแล้ว และควรกลับมาตรวจสอบว่ารั่วที่จุดใด

</aside>