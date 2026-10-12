# Go Fundamentals — คู่มือพื้นฐานภาษา Go (CS471 Day 3)

# 3. Securing a Go API

เอกสารนี้ต่อยอดจาก Day 2 โดยตรง โค้ดเริ่มต้นของวันนี้คือ `food-api` ที่ทำเสร็จเมื่อจบบทที่แล้ว หากยังทำไม่ครบ ให้ย้อนกลับไปทำถึงขั้นที่ 37 ของ Day 2 ก่อน

เป้าหมายของวันนี้คือปิดช่องโหว่หนึ่งข้อที่ค้างอยู่ตั้งแต่บทที่แล้ว แล้วเรียนรู้ว่าเรื่องความปลอดภัยวางอยู่ตรงไหนของ Clean Architecture

## 3.0 เตรียมตัวก่อนเริ่ม

### 3.0.1 ช่องโหว่ที่ค้างมาจากบทที่แล้ว

เปิดเซิร์ฟเวอร์ของบทที่แล้วขึ้นมา สร้างออเดอร์หนึ่งใบ แล้วเรียกใช้คำสั่งด้านล่าง

```bash
curl -i -X POST localhost:3000/api/orders/o1/cancel
```

```bash
HTTP/1.1 204 No Content
```

<aside>
🚨

คำสั่งนี้ไม่ได้แนบสิ่งใดที่บอกว่า **ผู้ใด** เป็นคนกด แต่ออเดอร์ถูกยกเลิกสำเร็จ

แปลว่าใครก็ตามที่เดาหมายเลขออเดอร์ได้ ยกเลิกออเดอร์ของผู้อื่นได้ทันที

ระบบที่เขียนมาทั้งบท ยังไม่รู้จักคำว่า *ผู้ใช้* เลยแม้แต่คำเดียว

</aside>

ทั้งวันนี้คือการปิดช่องโหว่นั้น โดยใช้กฎเดียวกับบทที่แล้วทุกประการ

<aside>
💡

**บทที่แล้ว** Postgres เป็นรายละเอียดของชั้น Repository ชั้น `domain` จึงไม่เคย import `pgx`

**บทนี้** JWT เป็นรายละเอียดของชั้น Delivery ชั้น `domain` จึงรู้จักเพียง `UserID` กับ `Role` และไม่เคยเห็นคำว่า Bearer

เป็นกฎข้อเดิม เพียงเปลี่ยนเรื่องที่นำมาใช้ และจบด้วยผลลัพธ์แบบเดียวกัน คือเปลี่ยนจาก JWT ไปเป็นระบบอื่นได้โดยแก้ไฟล์เดียว

</aside>

### 3.0.2 สิ่งที่ต้องมีในเครื่อง

| สิ่งที่ต้องมี | ตรวจด้วยคำสั่ง | หมายเหตุ |
| --- | --- | --- |
| โค้ดที่จบจาก Day 2 | `go test ./...` | ต้องผ่านทั้งหมดก่อนเริ่มบทนี้ |
| Docker Desktop | `docker info` | หรือใช้ Postgres ที่ติดตั้งไว้ในเครื่องเองก็ได้ เช่นเดียวกับบทที่แล้ว |
| `openssl` | `openssl version` | ใช้สร้างค่าสุ่มสำหรับ `JWT_SECRET` macOS และ Linux มีมาให้แล้ว |

### 3.0.3 วิธีอ่านเอกสารนี้

เนื้อหาแบ่งเป็น **17 ขั้น** ใช้กล่องสีแบบเดียวกับบทที่แล้ว คือกล่องสีเทาหมายถึงสร้างไฟล์ใหม่ และกล่องสีส้มหมายถึงแก้ไฟล์เดิม

<aside>
💡

ลำดับของบทนี้ไล่จากวงในออกวงนอกเช่นเดิม คือ `domain` ก่อน แล้วจึง `usecase` แล้วจึง Adapter และ `delivery` เป็นลำดับสุดท้าย

เหตุผลที่เรียงแบบนี้ เพราะกฎเรื่องสิทธิ์เป็นกฎทางธุรกิจ ไม่ใช่เรื่องของ HTTP จึงต้องเขียนให้เสร็จก่อนที่จะไปแตะ Token

</aside>

### 3.0.4 ผู้ใช้ Windows อ่านตรงนี้ก่อน

คำสั่งในเอกสารนี้เขียนด้วยไวยากรณ์ของ Shell ตระกูล Unix ซึ่งใช้ได้ทันทีบน macOS และ Linux ผู้ใช้ Windows มีสองทางเลือก

<aside>
✅

**ทางที่แนะนำ เปิด Git Bash แล้วทำงานในนั้นทั้งวัน**

Git for Windows ติดตั้ง Git Bash มาให้อยู่แล้ว ข้างในมีทั้ง `openssl` และ `curl` และรองรับ `$(...)` เหมือน macOS

เมื่อใช้ Git Bash คำสั่งทุกคำสั่งในเอกสารนี้ใช้ได้ตามที่เขียนไว้ ไม่ต้องแปลงอะไรเลย

ผู้ที่ติดตั้ง WSL ไว้แล้ว จะใช้ WSL แทนก็ได้ผลเหมือนกัน

</aside>

หากยืนยันจะใช้ PowerShell หรือ Command Prompt มีสามจุดที่ต้องระวัง

<aside>
🚨

**กับดักที่พบบ่อยที่สุด `curl` ใน Windows PowerShell ไม่ใช่ curl ตัวจริง**

ใน Windows PowerShell 5.1 ซึ่งเป็นรุ่นที่ติดมากับ Windows คำว่า `curl` เป็นชื่อเล่นของ `Invoke-WebRequest`

เมื่อพิมพ์ `curl -X POST ...` ตามเอกสาร จะได้ข้อความผิดพลาดที่อ่านไม่รู้เรื่อง เพราะคนละโปรแกรมกัน

**ทางแก้ที่ง่ายที่สุดคือพิมพ์ `curl.exe` แทน `curl`** ส่วนที่เหลือของคำสั่งใช้ได้ตามเดิม

Windows 10 รุ่น 1803 ขึ้นไปมี `curl.exe` ตัวจริงติดมาให้อยู่แล้ว

</aside>

ตารางเทียบคำสั่งที่ใช้บ่อยในเอกสารนี้

| สิ่งที่ต้องการ | Git Bash, macOS, Linux | Windows PowerShell |
| --- | --- | --- |
| เรียก API | `curl ...` | `curl.exe ...` |
| สร้างค่าสุ่ม 48 ไบต์ | `openssl rand -base64 48` | ดูคำสั่งเต็มในขั้นที่ 6 |
| คัดลอกไฟล์ | `cp .env.example .env` | `copy .env.example .env` |
| เก็บ Token ไว้ในตัวแปร | `TOKEN=$(... | jq -r .token)` | ดูคำสั่งเต็มด้านล่าง |
| ส่ง Token ไปกับ Request | `-H "Authorization: Bearer $TOKEN"` | `-H "Authorization: Bearer $TOKEN"` |

บน PowerShell การเก็บ Token ทำได้โดยไม่ต้องติดตั้ง `jq` เพราะ `Invoke-RestMethod` แปลง JSON ให้เองอยู่แล้ว

```powershell
$body = '{"email":"alice@example.com","password":"password123"}'
$TOKEN = (Invoke-RestMethod -Uri http://localhost:3000/api/auth/login -Method Post -ContentType 'application/json' -Body $body).token
$TOKEN
```

**คำสั่งวนซ้ำ ใช้ในขั้นที่ 12 ตอนทดสอบตัวจำกัดจำนวน Request**

ขั้นที่ 12 ใช้การวนยิงเจ็ดครั้งเพื่อดูว่าครั้งที่หกเป็นต้นไปถูกปฏิเสธ รูปแบบ `for ... do ... done` เป็นไวยากรณ์ของ Shell ตระกูล Unix **ใช้ใน PowerShell และ Command Prompt ไม่ได้**

```bash
for i in 1 2 3 4 5 6 7; do
  curl -s -o /dev/null -w "ครั้งที่ $i -> %{http_code}\n" \
    -X POST localhost:3000/api/auth/login -H 'Content-Type: application/json' \
    -d '{"email":"alice@example.com","password":"password123"}'
done
```

บน PowerShell ต้องเขียนใหม่เป็นแบบนี้

```powershell
$body = '{"email":"alice@example.com","password":"password123"}'
foreach ($i in 1..7) {
  try {
    $r = Invoke-WebRequest -Uri http://localhost:3000/api/auth/login -Method Post -ContentType 'application/json' -Body $body
    $code = $r.StatusCode
  } catch {
    $code = $_.Exception.Response.StatusCode.value__
  }
  Write-Host "ครั้งที่ $i -> $code"
}
```

<aside>
🚨

**`try` และ `catch` ในโค้ดด้านบนไม่ได้ใส่มาเผื่อ ๆ แต่จำเป็นจริง**

`Invoke-WebRequest` ถือว่า Status Code ระดับ 4xx และ 5xx เป็นข้อผิดพลาด แล้ว **โยน Exception ทิ้ง** แทนที่จะคืนค่ากลับมาเฉย ๆ

ถ้าไม่ดักไว้ การวนซ้ำจะหยุดทันทีที่เจอ 429 ซึ่งเป็นผลลัพธ์ที่เรากำลังอยากเห็นพอดี

ผู้ใช้ PowerShell 7 ขึ้นไป จะใส่ `-SkipHttpErrorCheck` แทนการดักด้วย `try` ก็ได้ แต่ตัวเลือกนี้ไม่มีใน Windows PowerShell 5.1

</aside>

<aside>
💡

ผู้สอนตรวจคำสั่งฝั่ง Unix ทั้งหมดด้วยการรันจริงแล้ว ส่วนคำสั่ง PowerShell เขียนขึ้นจากเอกสารของ .NET และ PowerShell แต่ไม่ได้รันทดสอบบนเครื่อง Windows

นี่คือเหตุผลที่แนะนำ Git Bash เป็นทางหลัก เพราะทุกคำสั่งในเอกสารนี้ผ่านการรันจริงมาแล้วในรูปแบบนั้น

</aside>

## 3.1 Authentication and Authorization

### 3.1.1 คำสองคำที่มักถูกใช้สลับกัน

| คำ | ตอบคำถามว่า |
| --- | --- |
| **Authentication** | *คุณคือใคร* พิสูจน์ตัวตนด้วยรหัสผ่านหรือ Token ผลลัพธ์คือรู้ว่าเป็นผู้ใช้คนใด |
| **Authorization** | *คุณทำสิ่งนี้ได้หรือไม่* ตัดสินจากตัวตนที่พิสูจน์แล้ว บวกกับกฎของระบบ |

<aside>
⚠️

ทั้งสองคำขึ้นต้นด้วย Auth เหมือนกัน จึงมักถูกเรียกรวมกันจนสับสน

วิธีจำที่ใช้ได้เสมอคือ `401 Unauthorized` แปลว่า *ยังไม่รู้ว่าคุณเป็นใคร* ส่วน `403 Forbidden` แปลว่า *รู้แล้วว่าคุณเป็นใคร แต่คุณทำสิ่งนี้ไม่ได้*

ชื่อของ `401` ในมาตรฐาน HTTP ใช้คำว่า Unauthorized ซึ่งสื่อความหมายผิด ที่ถูกควรเป็น Unauthenticated

</aside>

ในสถาปัตยกรรมของเรา สองเรื่องนี้อยู่คนละชั้นกันชัดเจน

| เรื่อง | อยู่ชั้นใด |
| --- | --- |
| แปลง Token เป็นตัวตน | `delivery` เพราะเป็นเรื่องของรูปแบบที่ส่งมาทาง HTTP |
| ตรวจรหัสผ่าน และออก Token | `usecase` โดยเรียกใช้ Adapter ที่ตนประกาศ Interface ไว้ |
| ตัดสินว่าใครแตะออเดอร์ใบนี้ได้ | `domain` เพราะเป็นกฎที่เป็นจริงเสมอ ไม่ว่าจะเรียกผ่าน HTTP หรือไม่ |

### 3.1.2 JWT ประกอบด้วยอะไร

JWT คือ **ข้อความธรรมดาสามท่อน ต่อกันด้วยจุด** ไม่มีการเข้ารหัสลับใด ๆ ทั้งสิ้น

```mermaid
flowchart LR
    T["JWT หนึ่งใบ"] --> H["ท่อนที่ 1 Header<br/>ใช้อัลกอริทึมใด"]
    T --> P["ท่อนที่ 2 Payload<br/>เป็นใคร สิทธิ์ใด หมดอายุเมื่อใด"]
    T --> S["ท่อนที่ 3 Signature<br/>ลายเซ็นที่คำนวณจากสองท่อนแรก"]
    H --> E1["Base64url อ่านกลับได้"]
    P --> E2["Base64url อ่านกลับได้"]
    S --> E3["HMAC-SHA256 ปลอมไม่ได้ถ้าไม่มีกุญแจ"]
```

ตลอดหัวข้อนี้จะใช้ **โทเคนจริงหนึ่งใบ** เป็นตัวอย่าง ทุกคำสั่งด้านล่างรันตามได้เลยและจะได้ผลเหมือนกันทุกเครื่อง

```bash
SECRET='cs471-demo-secret-0123456789-abcdefghij'

TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1MSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MjAwMDAwMDAwMH0.5yRKOzmonQCqgPsFHx1H2rFrtwpvaTFjSHZ8SusOpco
```

แยกเป็นสามท่อนตามจุดจะเห็นชัดขึ้น

```bash
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9                              <- Header    36 ตัวอักษร
eyJzdWIiOiJ1MSIsInJvbGUiOiJjdXN0b21lciIsImV4cCI6MjAwMDAwMDAwMH0   <- Payload   63 ตัวอักษร
5yRKOzmonQCqgPsFHx1H2rFrtwpvaTFjSHZ8SusOpco                       <- Signature 43 ตัวอักษร
```

**ท่อนที่สามยาว 43 ตัวอักษรเสมอ ไม่ว่า Payload จะยาวแค่ไหน** เพราะ SHA256 ให้ผลลัพธ์ 256 บิตคงที่ ซึ่งเมื่อเข้ารหัส Base64url แบบไม่มี Padding จะได้ 43 ตัวอักษรพอดี

**ถอดสองท่อนแรกกลับเป็นข้อความ**

Base64url ตัดเครื่องหมาย `=` ท้ายทิ้ง ก่อนถอดจึงต้องเติมกลับให้ความยาวหารด้วย 4 ลงตัว ฟังก์ชันเล็ก ๆ ด้านล่างทำให้ทั้งหมดนั้น

```bash
d() {
  local s=$(echo "$1" | tr '_-' '/+')
  local m=$(( ${#s} % 4 ))
  [ $m -ne 0 ] && s="$s$(printf '=%.0s' $(seq $((4-m))))"
  echo "$s" | openssl base64 -d -A
}

d $(echo $TOKEN | cut -d. -f1)   # Header
d $(echo $TOKEN | cut -d. -f2)   # Payload
```

```bash
{"alg":"HS256","typ":"JWT"}
{"sub":"u1","role":"customer","exp":2000000000}
```

<aside>
🚨

**สังเกตให้ดี ไม่ได้ใส่กุญแจอะไรเลย แต่อ่านออกทั้งหมด**

Base64 คือการ *เปลี่ยนรูป* ไม่ใช่การ *เข้ารหัสลับ* ใครก็ตามที่ได้ Token ไป เปิดอ่าน Payload ได้ทันที

จึงห้ามใส่รหัสผ่าน เลขบัตรประชาชน หรือข้อมูลลับใด ๆ ลงใน Payload เด็ดขาด

สิ่งที่ลายเซ็นรับประกันคือ *ข้อมูลนี้ไม่ถูกแก้ระหว่างทาง* ไม่ใช่ *ข้อมูลนี้เป็นความลับ*

</aside>

**ท่อนที่สามคำนวณมาอย่างไร**

ลายเซ็นคือผลของฟังก์ชัน HMAC-SHA256 ที่รับสองอย่างเข้าไป คือ *ข้อความที่จะเซ็น* กับ *กุญแจลับ*

```bash
signature = Base64url( HMAC-SHA256( Base64url(header) + "." + Base64url(payload), SECRET ) )
```

```mermaid
flowchart TD
    A["Base64url ของ Header"] --> C["ต่อด้วยจุด กลายเป็นข้อความที่จะเซ็น"]
    B["Base64url ของ Payload"] --> C
    K["SECRET กุญแจลับ อยู่แต่ในเซิร์ฟเวอร์"] --> D["HMAC-SHA256"]
    C --> D
    D --> E["ได้ 256 บิต แปลงเป็น Base64url 43 ตัวอักษร"]
    E --> F["นำไปต่อท้ายเป็นท่อนที่ 3"]
```

**พิสูจน์ด้วยมือ** คำนวณลายเซ็นเองแล้วเทียบกับท่อนที่สามในโทเคน

```bash
printf '%s' "$(echo $TOKEN | cut -d. -f1).$(echo $TOKEN | cut -d. -f2)" \
  | openssl dgst -sha256 -hmac "$SECRET" -binary \
  | openssl base64 -A | tr '+/' '-_' | tr -d '='
```

```bash
5yRKOzmonQCqgPsFHx1H2rFrtwpvaTFjSHZ8SusOpco
```

<aside>
✅

ค่าที่คำนวณเองตรงกับท่อนที่สามในโทเคนทุกตัวอักษร

แปลว่า JWT **ไม่มีอะไรลึกลับเลย** เป็นเพียง Base64 สองท่อน บวกผลของ HMAC หนึ่งท่อน

ทั้งใบสร้างได้ด้วย `openssl` อย่างเดียว ไม่ต้องมีไลบรารีใด และโทเคนที่สร้างแบบนี้ เซิร์ฟเวอร์ที่เราจะเขียนวันนี้ก็ยอมรับ

</aside>

**แล้วทำไมปลอมไม่ได้**

ลองสวมบทผู้โจมตี แก้ `role` จาก `customer` เป็น `staff` แล้วประกอบโทเคนใหม่ ปัญหาอยู่ที่ท่อนที่สาม เพราะผู้โจมตีไม่มี `SECRET` จึงได้แค่เดา

```bash
FAKE=$(printf '%s' '{"sub":"u1","role":"staff","exp":2000000000}' \
  | openssl base64 -A | tr '+/' '-_' | tr -d '=')

# ผู้โจมตีเดากุญแจ จึงได้ลายเซ็นคนละตัวกับที่ควรเป็น
printf '%s' "$(echo $TOKEN | cut -d. -f1).$FAKE" \
  | openssl dgst -sha256 -hmac 'wrong-guess-secret' -binary \
  | openssl base64 -A | tr '+/' '-_' | tr -d '='
```

```bash
ลายเซ็นที่ผู้โจมตีเดาได้ : dfQmjBxGQlts8R_MDoHVOJCc5873VvBAp9fdeAC9t5M
ลายเซ็นที่ถูกต้องจริง   : wAWzgzhrgXnNg_zH528uYS6RmrEDTeAvH6ivvO9Po6U
```

เซิร์ฟเวอร์คำนวณลายเซ็นที่ถูกต้องขึ้นมาเองจาก Payload ที่ได้รับ แล้วเทียบกับท่อนที่สามที่แนบมา ไม่ตรงก็ปฏิเสธ จบ

- 🧠 เจาะลึก ทำไมการเดาจึงไม่มีทางสำเร็จ
    
    ```bash
    เปลี่ยนข้อความที่เซ็นเพียงตัวอักษรเดียว จาก sub=u1 เป็น sub=u2
    โดยใช้กุญแจตัวเดิม ผลลัพธ์ที่ได้คือ
    
      u1 -> e7244a3b39a89d00aa80fb051f1d47dab16bb70a6f69316348767c4aeb0ea5ca
      u2 -> 4a8a5a27920c1dbaf030a8b92f19d29e30b9d9b59d2917f94de830bafe0494fb
    
    ต่างกัน 109 บิต จาก 256 บิต คิดเป็น 43 เปอร์เซ็นต์
    และไม่มีไบต์ใดเลยใน 32 ไบต์ที่บังเอิญเหมือนกัน
    
    คุณสมบัตินี้เรียกว่า Avalanche Effect คือเปลี่ยน input นิดเดียว
    output ต้องเปลี่ยนไปราวครึ่งหนึ่งแบบคาดเดาไม่ได้
    ผู้โจมตีจึงไล่เดาทีละบิตแบบเข้าใกล้คำตอบไปเรื่อย ๆ ไม่ได้
    เพราะ "ใกล้เคียง" ไม่ได้ให้สัญญาณอะไรเลย
    
    ทางเดียวที่เหลือคือไล่ลองกุญแจทุกความเป็นไปได้
    ซึ่งเป็นเหตุผลที่ JWT_SECRET ต้องยาวและสุ่มจริง
    กุญแจที่เดาได้ เช่น secret หรือ 12345678 ทำให้ทุกอย่างข้างบนนี้พังทันที
    ```
    

**เซิร์ฟเวอร์ตรวจอย่างไร ทีละขั้น**

```mermaid
sequenceDiagram
    participant C as ผู้ใช้
    participant S as เซิร์ฟเวอร์
    C->>S: POST /auth/login อีเมลกับรหัสผ่าน
    Note over S: ตรวจรหัสผ่านผ่าน แล้วเซ็นโทเคนด้วย SECRET
    S-->>C: ส่ง JWT กลับไป
    Note over C: เก็บไว้ แล้วแนบไปทุก Request<br/>Authorization: Bearer ...
    C->>S: GET /orders/o1 พร้อมโทเคน
    Note over S: 1 แยกสามท่อน<br/>2 ตรวจว่า alg เป็น HS256 จริง<br/>3 คำนวณลายเซ็นใหม่จากสองท่อนแรก<br/>4 เทียบกับท่อนที่สาม<br/>5 ตรวจ exp ว่ายังไม่หมดอายุ
    S-->>C: ผ่านครบห้าข้อจึงตอบข้อมูล ไม่เช่นนั้น 401
```

<aside>
💡

สังเกตว่าทั้งห้าขั้นตอนนั้น **เซิร์ฟเวอร์ไม่ต้องเปิดฐานข้อมูลเลยแม้แต่ครั้งเดียว**

ทุกอย่างที่ต้องใช้ตัดสินอยู่ในตัวโทเคนเอง บวกกับกุญแจที่เซิร์ฟเวอร์มีอยู่แล้ว

นี่คือความหมายของคำว่า **Stateless** และเป็นทั้งข้อดีและข้อเสียของ JWT ที่จะเทียบกันในหัวข้อ 3.11

</aside>

- 🧠 เจาะลึก HS256 กับ RS256 ต่างกันอย่างไร
    
    ```bash
    HS256 ที่ใช้ในบทนี้เป็นแบบ สมมาตร คือกุญแจเซ็นกับกุญแจตรวจเป็นดอกเดียวกัน
    ใครตรวจได้ก็เซ็นได้ จึงเหมาะกับกรณีที่ผู้เซ็นและผู้ตรวจเป็นระบบเดียวกัน
    อย่างที่เราทำวันนี้ คือ Fiber ตัวเดียวทั้งออกและตรวจโทเคน
    
    RS256 เป็นแบบ อสมมาตร มีกุญแจสองดอกคู่กัน
      กุญแจส่วนตัว  ใช้เซ็น    เก็บไว้ที่ผู้ออกโทเคนเท่านั้น
      กุญแจสาธารณะ ใช้ตรวจ    แจกให้ใครก็ได้
    
    ใช้เมื่อผู้ตรวจกับผู้ออกไม่ใช่คนเดียวกัน เช่นมี Auth Service กลางออกโทเคน
    แล้วมี Service อีกสิบตัวที่ต้องตรวจ ถ้าใช้ HS256 ต้องแจกกุญแจลับให้ทั้งสิบตัว
    ซึ่งแปลว่าทั้งสิบตัวปลอมโทเคนได้ RS256 แจกเฉพาะกุญแจสาธารณะ ตรวจได้แต่เซ็นไม่ได้
    
    โครงสร้างที่เขียนวันนี้เปลี่ยนไปใช้ RS256 ได้โดยแก้ security/jwt.go ไฟล์เดียว
    เพราะชั้นอื่นรู้จักแค่ Interface ชื่อ TokenIssuer
    ```
    

### 3.1.3 Token กับ Session ต่างกันตรงไหน

| Session ในหน่วยความจำ | Token |
| --- | --- |
| เซิร์ฟเวอร์ต้องจำว่าใครล็อกอินอยู่บ้าง | เซิร์ฟเวอร์ไม่ต้องจำอะไรเลย ทุกอย่างอยู่ใน Token |
| เป็น **Stateful** เพิ่มสำเนาแล้วพัง เว้นแต่ย้าย Session ไป Redis | เป็น **Stateless** สำเนาใดก็ตรวจ Token ได้เหมือนกัน |
| ยกเลิกได้ทันที เพียงลบ Session ทิ้ง | ยกเลิกก่อนหมดอายุไม่ได้ ต้องรอหมดอายุเอง |

<aside>
💡

ข้อได้เปรียบและข้อเสียเปรียบสลับกันพอดี Token แลกความสามารถในการยกเลิกทันที กับการที่เซิร์ฟเวอร์ไม่ต้องเก็บสถานะ

นี่คือหัวข้อ 2.5.1 Stateless และ Stateful ของบทที่แล้วโดยตรง เพียงเปลี่ยนจากเรื่องข้อมูลมาเป็นเรื่องตัวตน

ด้วยเหตุนี้ Token ในระบบจริงจึงมีอายุสั้น และมักใช้คู่กับ Refresh Token

วันนี้จะเริ่มจาก Token ก่อน เพราะเห็นภาพง่ายกว่า แล้วใน **หัวข้อ 3.11** จะสร้างแบบ Session ขึ้นมาอีกตัวหนึ่ง

เสียบเข้าที่เดิมโดยไม่แก้ชั้น Domain และ Use Case เลย แล้วค่อยเทียบข้อดีข้อเสียกันด้วยผลที่วัดได้จริง

</aside>

## 3.2 ชั้นที่ 1 — ตัวตนและกฎเรื่องสิทธิ์อยู่ใน Domain

เริ่มจากวงในสุดเช่นเดิม สิ่งแรกที่ระบบต้องรู้จักคือ *ผู้ใช้หนึ่งคน*

#### ขั้นที่ 1 — ชั้น Domain — User, Role และ Error ของเรื่องสิทธิ์

<aside>
✏️

**ขั้นนี้มีทั้งไฟล์ใหม่และไฟล์เดิม เป็นการสอนคำศัพท์ใหม่ให้ชั้นในสุด**

**1. สร้างไฟล์** `domain/user.go` แล้ววางโค้ดด้านล่างลงไป

**2. เปิดไฟล์** `domain/errors.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 13 แล้วเพิ่ม Error ตามรายการด้านล่าง

**3. รัน** `go build ./...`

ขั้นนี้มี **2 ไฟล์** ทำให้ครบทั้งคู่ก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `domain/user.go` สร้างใหม่**

```go
package domain

// User คือผู้ใช้หนึ่งคนในระบบ
//
// สังเกตว่าไม่มีคำว่า JWT, Token หรือ Bearer อยู่ในไฟล์นี้เลย
// ชั้นในสุดรู้จักเพียงว่า "ผู้ใช้คนนี้คือใคร และมีสิทธิ์ระดับใด"
// ส่วนวิธีพิสูจน์ตัวตนเป็นเรื่องของชั้นนอก
type User struct {
	ID           string
	Email        string
	PasswordHash string // เก็บเฉพาะค่าที่ผ่านการ Hash แล้ว ไม่เคยเก็บรหัสผ่านจริง
	Role         Role
}

// Role ประกาศเป็น Type ของตัวเอง ไม่ใช่ string ธรรมดา
// ด้วยเหตุผลเดียวกับ Status ในบทที่แล้ว คือให้คอมไพเลอร์ช่วยตรวจการพิมพ์ผิด
type Role string

const (
	Customer Role = "customer" // ลูกค้า จัดการได้เฉพาะออเดอร์ของตนเอง
	Staff    Role = "staff"    // พนักงานร้าน เห็นและจัดการได้ทุกออเดอร์
)

func (u User) IsStaff() bool { return u.Role == Staff }

// Valid ใช้ตรวจว่าค่าที่รับเข้ามาเป็น Role ที่มีอยู่จริงหรือไม่
func (r Role) Valid() bool {
	switch r {
	case Customer, Staff:
		return true
	}
	return false
}
```

- ไฟล์นี้ไม่ import สิ่งใดเลย เช่นเดียวกับ `menu.go` และ `order.go` ในบทที่แล้ว
- `Role` ประกาศเป็น Type ของตัวเองด้วยเหตุผลเดียวกับ `Status` คือให้คอมไพเลอร์ช่วยจับการพิมพ์ผิด
- ฟีลด์ชื่อ `PasswordHash` ไม่ใช่ `Password` เพราะระบบนี้ไม่เคยเก็บรหัสผ่านจริงไว้ที่ใดเลย

**ไฟล์ที่ 2 — `domain/errors.go` สิ่งที่แก้ไข**

เพิ่มบรรทัดด้านล่างนี้ต่อท้าย ภายในบล็อก `var` เดิม

```go
// เพิ่มในบทนี้ ทั้งหมดยังคงเป็น Sentinel Error ธรรมดา
// ชั้นในคืนค่าออกไป ส่วนชั้นนอกเป็นผู้แปลงเป็น Status Code เช่นเดิม
ErrUnauthorized       = errors.New("unauthorized")              // ยังไม่ได้พิสูจน์ตัวตน
ErrForbidden          = errors.New("forbidden")                 // พิสูจน์ตัวตนแล้ว แต่ไม่มีสิทธิ์
ErrInvalidCredentials = errors.New("invalid email or password") // อีเมลหรือรหัสผ่านไม่ถูกต้อง
ErrEmailTaken         = errors.New("email already registered")  // สมัครซ้ำ
ErrWeakPassword       = errors.New("password must be at least 8 characters")
ErrTooManyRequests    = errors.New("too many requests") // เรียกถี่เกินกำหนด
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `domain/errors.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    
    	// เพิ่มในบทนี้ ทั้งหมดยังคงเป็น Sentinel Error ธรรมดา
    	// ชั้นในคืนค่าออกไป ส่วนชั้นนอกเป็นผู้แปลงเป็น Status Code เช่นเดิม
    	ErrUnauthorized       = errors.New("unauthorized")              // ยังไม่ได้พิสูจน์ตัวตน
    	ErrForbidden          = errors.New("forbidden")                 // พิสูจน์ตัวตนแล้ว แต่ไม่มีสิทธิ์
    	ErrInvalidCredentials = errors.New("invalid email or password") // อีเมลหรือรหัสผ่านไม่ถูกต้อง
    	ErrEmailTaken         = errors.New("email already registered")  // สมัครซ้ำ
    	ErrWeakPassword       = errors.New("password must be at least 8 characters")
    	ErrTooManyRequests    = errors.New("too many requests")                    // เรียกถี่เกินกำหนด
    	ErrRevokeUnsupported  = errors.New("logout is not supported in this mode") // Token ที่ยกเลิกไม่ได้
    	ErrIDTaken            = errors.New("id already used")                      // ชนกับ ID ที่มีอยู่แล้ว
    )
    ```
    
- ยังคงเป็นรูปแบบเดิมทุกประการ คือชั้นในประกาศและคืนค่า ส่วนชั้นนอกเป็นผู้ตัดสินว่าตัวใดควรเป็น Status Code ใด
- `ErrInvalidCredentials` ใช้ข้อความเดียวสำหรับทั้งกรณีอีเมลผิดและรหัสผ่านผิด ด้วยเหตุผลที่จะอธิบายที่ **ขั้นที่ 3**

#### ขั้นที่ 2 — ชั้น Domain — กฎเรื่องสิทธิ์ พร้อม Test ที่ไม่ต้องมี Token

<aside>
✏️

**ขั้นนี้เขียนกฎแล้วพิสูจน์กฎในขั้นเดียวกัน**

**1. เปิดไฟล์** `domain/order.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 12 แล้วแก้ตามรายการด้านล่าง

**2. สร้างไฟล์** `domain/permission_test.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go test ./domain/ -v`

ขั้นนี้มี **2 ไฟล์** ทำให้ครบทั้งคู่ก่อนรันคำสั่ง กฎที่เพิ่งเขียนจะถูกพิสูจน์ทันทีในคำสั่งเดียว

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `domain/order.go` สิ่งที่แก้ไข**

- เพิ่มฟีลด์ `UserID` เข้าไปใน `Order` เพื่อบันทึกว่าออเดอร์ใบนี้เป็นของผู้ใด
- เพิ่ม 4 Method ที่ตอบคำถามเรื่องสิทธิ์ ต่อท้ายไฟล์

เพิ่มฟีลด์เข้าไปใน Struct `Order` ให้เป็นดังนี้

```go
type Order struct {
	ID     string
	UserID string // เจ้าของออเดอร์ ใช้ตัดสินว่าใครมีสิทธิ์แตะออเดอร์ใบนี้
	Lines  []OrderLine
	Status Status
}
```

จากนั้นเพิ่มโค้ดด้านล่างนี้ต่อท้ายไฟล์

```go
// OwnedBy ตอบว่าออเดอร์ใบนี้เป็นของผู้ใช้คนนั้นหรือไม่
func (o Order) OwnedBy(userID string) bool { return o.UserID == userID }

// CanBeCancelledBy ต่อยอดจาก CanCancel ของบทที่แล้ว
// โดยเพิ่มคำถามว่า "ผู้ที่กดยกเลิก มีสิทธิ์กับออเดอร์ใบนี้หรือไม่"
func (o Order) CanBeCancelledBy(u User) bool {
	if !o.CanCancel() {
		return false
	}
	return u.IsStaff() || o.OwnedBy(u.ID)
}

// CanBeViewedBy ลูกค้าเห็นเฉพาะออเดอร์ของตนเอง ส่วนพนักงานเห็นทุกใบ
func (o Order) CanBeViewedBy(u User) bool {
	return u.IsStaff() || o.OwnedBy(u.ID)
}

// การเดินสถานะเป็นงานของร้าน ลูกค้าเปลี่ยนสถานะออเดอร์เองไม่ได้
func (o Order) CanChangeStatusBy(u User) bool { return u.IsStaff() }
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `domain/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    	UserID string // เจ้าของออเดอร์ ใช้ตัดสินว่าใครมีสิทธิ์แตะออเดอร์ใบนี้
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
    
    // ─────────────── กฎเรื่องสิทธิ์ อยู่ชั้นเดียวกับกฎธุรกิจอื่น ───────────────
    
    // OwnedBy ตอบว่าออเดอร์ใบนี้เป็นของผู้ใช้คนนั้นหรือไม่
    func (o Order) OwnedBy(userID string) bool { return o.UserID == userID }
    
    // CanBeCancelledBy ต่อยอดจาก CanCancel ของบทที่แล้ว
    // โดยเพิ่มคำถามว่า "ผู้ที่กดยกเลิก มีสิทธิ์กับออเดอร์ใบนี้หรือไม่"
    //
    // เรียงลำดับการตรวจให้ถูก คือตรวจสถานะก่อน แล้วจึงตรวจสิทธิ์
    // เพราะออเดอร์ที่กำลังทำอาหารอยู่ ยกเลิกไม่ได้ไม่ว่าผู้ใดจะเป็นคนกด
    func (o Order) CanBeCancelledBy(u User) bool {
    	if !o.CanCancel() {
    		return false
    	}
    	return u.IsStaff() || o.OwnedBy(u.ID)
    }
    
    // CanBeViewedBy ลูกค้าเห็นเฉพาะออเดอร์ของตนเอง ส่วนพนักงานเห็นทุกใบ
    func (o Order) CanBeViewedBy(u User) bool {
    	return u.IsStaff() || o.OwnedBy(u.ID)
    }
    
    // การเดินสถานะเป็นงานของร้าน ลูกค้าเปลี่ยนสถานะออเดอร์เองไม่ได้
    func (o Order) CanChangeStatusBy(u User) bool { return u.IsStaff() }
    ```
    
- `CanBeCancelledBy` **เรียกใช้** `CanCancel` ของเดิม ไม่ได้เขียนกฎสถานะขึ้นใหม่ กฎเดิมจึงยังมีที่อยู่เพียงแห่งเดียว
- ลำดับการตรวจสำคัญ ตรวจสถานะก่อน แล้วจึงตรวจสิทธิ์ เพราะออเดอร์ที่กำลังทำอาหารอยู่ ยกเลิกไม่ได้ไม่ว่าผู้ใดจะเป็นคนกด แม้แต่พนักงาน
- ทั้งสี่ Method ไม่รู้จัก HTTP ไม่รู้จัก Token และไม่รู้จักฐานข้อมูล จึงทดสอบได้ทันทีในขั้นถัดไป

**ไฟล์ที่ 2 — `domain/permission_test.go` สร้างใหม่**

```go
package domain_test

import (
	"testing"

	"food-api/domain"
)

// กฎเรื่องสิทธิ์อยู่ในชั้น domain จึงทดสอบได้โดยไม่ต้องมี Token
// ไม่ต้องมีฐานข้อมูล และไม่ต้องเปิดเซิร์ฟเวอร์ เหมือนกฎธุรกิจอื่นในบทที่แล้ว
var (
	owner  = domain.User{ID: "u1", Role: domain.Customer}
	other  = domain.User{ID: "u2", Role: domain.Customer}
	worker = domain.User{ID: "u9", Role: domain.Staff}
)

func TestCanBeCancelledBy(t *testing.T) {
	cases := []struct {
		name   string
		order  domain.Order
		actor  domain.User
		expect bool
	}{
		{"เจ้าของยกเลิกออเดอร์ที่ยัง pending ได้", domain.Order{UserID: "u1", Status: domain.Pending}, owner, true},
		{"คนอื่นยกเลิกออเดอร์ของเราไม่ได้", domain.Order{UserID: "u1", Status: domain.Pending}, other, false},
		{"พนักงานยกเลิกออเดอร์ของลูกค้าได้", domain.Order{UserID: "u1", Status: domain.Pending}, worker, true},
		{"เริ่มทำอาหารแล้ว เจ้าของก็ยกเลิกไม่ได้", domain.Order{UserID: "u1", Status: domain.Cooking}, owner, false},
		{"เริ่มทำอาหารแล้ว พนักงานก็ยกเลิกไม่ได้", domain.Order{UserID: "u1", Status: domain.Cooking}, worker, false},
	}

	for _, c := range cases {
		if got := c.order.CanBeCancelledBy(c.actor); got != c.expect {
			t.Errorf("%s: ได้ %v ต้องการ %v", c.name, got, c.expect)
		}
	}
}

func TestCanBeViewedBy(t *testing.T) {
	o := domain.Order{UserID: "u1", Status: domain.Pending}

	if !o.CanBeViewedBy(owner) {
		t.Error("เจ้าของต้องเห็นออเดอร์ของตนเอง")
	}
	if o.CanBeViewedBy(other) {
		t.Error("ลูกค้าคนอื่นต้องไม่เห็นออเดอร์ใบนี้")
	}
	if !o.CanBeViewedBy(worker) {
		t.Error("พนักงานต้องเห็นได้ทุกใบ")
	}
}

func TestCanChangeStatusBy(t *testing.T) {
	o := domain.Order{UserID: "u1", Status: domain.Pending}

	if o.CanChangeStatusBy(owner) {
		t.Error("ลูกค้าต้องเปลี่ยนสถานะเองไม่ได้ แม้เป็นออเดอร์ของตน")
	}
	if !o.CanChangeStatusBy(worker) {
		t.Error("พนักงานต้องเปลี่ยนสถานะได้")
	}
}
```

<aside>
💡

นี่คือผลตอบแทนของการวางกฎไว้ในชั้น `domain` ตั้งแต่บทที่แล้ว

กฎเรื่องสิทธิ์ทั้งหมดพิสูจน์ได้ในเสี้ยววินาที โดยไม่ต้องสมัครสมาชิก ไม่ต้องขอ Token ไม่ต้องเปิดฐานข้อมูล และไม่ต้องส่ง HTTP Request

หากเขียนกฎเหล่านี้ไว้ใน Handler จะทดสอบแบบนี้ไม่ได้เลย

</aside>

**ผลลัพธ์**

```bash
go test ./domain/ -run 'TestCanBeCancelledBy|TestCanBeViewedBy|TestCanChangeStatusBy' -v
```

```bash
=== RUN   TestCanBeCancelledBy
--- PASS: TestCanBeCancelledBy (0.00s)
=== RUN   TestCanBeViewedBy
--- PASS: TestCanBeViewedBy (0.00s)
=== RUN   TestCanChangeStatusBy
--- PASS: TestCanChangeStatusBy (0.00s)
PASS
ok  	food-api/domain	1.165s
```

## 3.3 ชั้นที่ 2 — Use Case ประกาศสิ่งที่ตนต้องการ

ชั้นนี้ต้อง Hash รหัสผ่านและออก Token แต่ **ต้องไม่รู้จักชื่อ bcrypt และ JWT** จึงใช้วิธีเดียวกับ `OrderRepository` ในบทที่แล้ว คือประกาศ Interface ไว้ฝั่งผู้ใช้งาน

#### ขั้นที่ 3 — ชั้น Use Case — ประกาศ Interface สามตัว แล้วเขียน Register และ Login

<aside>
📄

**1. สร้างไฟล์** `usecase/auth.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package usecase

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"food-api/domain"
)

// Interface สามตัวนี้ประกาศไว้ที่ชั้น Use Case ซึ่งเป็นผู้ใช้งาน
// ไม่ใช่ที่ชั้นนอกซึ่งเป็นผู้ให้บริการ หลักการเดียวกับ OrderRepository ในบทที่แล้ว
//
// ผลคือชั้นนี้ไม่รู้จัก bcrypt และไม่รู้จัก JWT เลยแม้แต่ชื่อ

type UserRepository interface {
	CreateUser(u domain.User) error
	FindUserByEmail(email string) (domain.User, error)
}

// PasswordHasher ซ่อนวิธีเข้ารหัสไว้ข้างหลัง จะเป็น bcrypt หรือ argon2 ก็ได้
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) error
}

// TokenIssuer ซ่อนรูปแบบของ Token ไว้ข้างหลัง จะเป็น JWT หรือแบบอื่นก็ได้
type TokenIssuer interface {
	Issue(u domain.User) (string, error)
	Parse(token string) (domain.User, error)
}

type AuthUseCase struct {
	repo   UserRepository
	hasher PasswordHasher
	tokens TokenIssuer
	seq    *int64

	// dummyHash คือ Hash ของรหัสผ่านที่ไม่มีใครใช้จริง
	// มีไว้เผาเวลาให้เท่ากัน ตอนที่ค้นหาอีเมลแล้วไม่พบ ดูคำอธิบายใน Login
	dummyHash string
}

func NewAuthUseCase(r UserRepository, h PasswordHasher, t TokenIssuer) AuthUseCase {
	var n int64
	dummy, _ := h.Hash("a-password-that-nobody-will-ever-use")
	return AuthUseCase{repo: r, hasher: h, tokens: t, seq: &n, dummyHash: dummy}
}

func (a AuthUseCase) nextID() string {
	return fmt.Sprintf("u%d", atomic.AddInt64(a.seq, 1))
}

// Register คือทางสมัครสาธารณะ ผลลัพธ์เป็นลูกค้าเสมอ
//
// สังเกตว่าไม่มีพารามิเตอร์ role ให้ส่งเข้ามา ผู้สมัครจึงเลือกสิทธิ์ของตนเองไม่ได้
// แม้ชั้นนอกจะเผลอส่งค่ามาก็ตาม นี่คือการปิดช่องโหว่ Mass Assignment ที่ต้นทาง
func (a AuthUseCase) Register(email, password string) (domain.User, error) {
	return a.create(email, password, domain.Customer)
}

// CreateStaff ไม่มี Endpoint สาธารณะเส้นใดชี้มาที่นี่
// ใช้ตอนตั้งระบบเท่านั้น บัญชีพนักงานคนแรกจึงมาจากผู้ดูแลระบบ ไม่ใช่จากการสมัครเอง
func (a AuthUseCase) CreateStaff(email, password string) (domain.User, error) {
	return a.create(email, password, domain.Staff)
}

func (a AuthUseCase) create(email, password string, role domain.Role) (domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || !strings.Contains(email, "@") {
		return domain.User{}, domain.ErrInvalidCredentials
	}
	// ความยาวรหัสผ่านเป็นกฎของระบบ จึงตรวจที่ชั้นนี้ ไม่ใช่ที่ Handler
	if len(password) < 8 {
		return domain.User{}, domain.ErrWeakPassword
	}
	if !role.Valid() {
		role = domain.Customer
	}

	if _, err := a.repo.FindUserByEmail(email); err == nil {
		return domain.User{}, domain.ErrEmailTaken
	}

	hash, err := a.hasher.Hash(password)
	if err != nil {
		return domain.User{}, fmt.Errorf("register: %w", err)
	}

	// ตัวนับ ID เริ่มจากศูนย์ใหม่ทุกครั้งที่โปรแกรมรีสตาร์ท แต่ข้อมูลในฐานข้อมูลไม่ได้หายไปด้วย
	// หมายเลขที่ได้จึงอาจชนกับผู้ใช้ที่มีอยู่แล้ว วิธีแก้คือขยับไปหมายเลขถัดไปจนกว่าจะว่าง
	for i := 0; i < 100; i++ {
		u := domain.User{ID: a.nextID(), Email: email, PasswordHash: hash, Role: role}

		err := a.repo.CreateUser(u)
		if err == nil {
			return u, nil
		}
		// ชนเฉพาะเรื่อง ID เท่านั้นที่ลองใหม่ได้ ปัญหาอื่นให้คืนออกไปทันที
		if !errors.Is(err, domain.ErrIDTaken) {
			return domain.User{}, fmt.Errorf("register: %w", err)
		}
	}
	return domain.User{}, fmt.Errorf("register: หาหมายเลขผู้ใช้ที่ว่างไม่ได้")
}

func (a AuthUseCase) Login(email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	u, err := a.repo.FindUserByEmail(email)
	if err != nil {
		// ตอบข้อความเดียวกันทั้งกรณีไม่มีอีเมลนี้ และกรณีรหัสผ่านผิด
		// มิฉะนั้นผู้โจมตีจะใช้ API นี้ไล่เดาว่าอีเมลใดมีอยู่ในระบบบ้าง
		//
		// ข้อความเหมือนกันอย่างเดียวยังไม่พอ เพราะ bcrypt ใช้เวลาราว 50 มิลลิวินาที
		// ถ้าเส้นทางนี้คืนค่าทันทีโดยไม่ตรวจอะไรเลย ผู้โจมตีจับเวลาก็แยกออกทันที
		// จึงต้องเสียเวลาให้เท่ากัน ด้วยการตรวจกับ Hash ปลอมที่ไม่มีวันตรงกับอะไร
		_ = a.hasher.Compare(a.dummyHash, password)
		return "", domain.ErrInvalidCredentials
	}

	if err := a.hasher.Compare(u.PasswordHash, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	return a.tokens.Issue(u)
}

// Identify แปลง Token กลับเป็นตัวตน ใช้โดย Middleware ที่ชั้นนอก
func (a AuthUseCase) Identify(token string) (domain.User, error) {
	u, err := a.tokens.Parse(token)
	if err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}
	return u, nil
}
```

- Interface ทั้งสามตัวประกาศอยู่ใน Package `usecase` ซึ่งเป็นผู้ใช้งาน ไม่ใช่ใน Package ของผู้ให้บริการ เป็น Dependency Inversion ข้อเดิม
- ค้นหาคำว่า `bcrypt` และ `jwt` ในไฟล์นี้จะไม่พบเลย ชั้นนี้รู้เพียงว่ามีสิ่งที่ *Hash ได้* และสิ่งที่ *ออก Token ได้*
- กฎความยาวรหัสผ่านอย่างน้อย 8 ตัวอักษร อยู่ที่ชั้นนี้ ไม่ใช่ที่ Handler เพราะเป็นกฎของระบบ ไม่ใช่กฎของ HTTP

<aside>
🚨

สังเกต `Login` ให้ดี ทั้งกรณี **ไม่มีอีเมลนี้ในระบบ** และกรณี **รหัสผ่านผิด** คืน `ErrInvalidCredentials` ตัวเดียวกัน และข้อความเดียวกัน

หากแยกข้อความเป็น *ไม่พบอีเมลนี้* กับ *รหัสผ่านไม่ถูกต้อง* ผู้โจมตีจะใช้หน้าเข้าสู่ระบบไล่ตรวจได้ว่าอีเมลใดมีอยู่ในระบบบ้าง

ช่องโหว่ลักษณะนี้เรียกว่า User Enumeration และพบได้บ่อยมากในระบบที่เขียนขึ้นเอง

</aside>

#### ขั้นที่ 4 — ชั้น Use Case — ส่งตัวตนเข้าไปให้ทุก Method แล้วแก้ Test ตาม

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `usecase/order.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 7

**2. เปิดไฟล์** `usecase/order_test.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 27

**3. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**4. รัน** `go test ./domain/ ./usecase/`

ขั้นนี้มี **2 ไฟล์** แก้ให้ครบทั้งคู่ก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `usecase/order.go` สิ่งที่แก้ไข**

- ทุก Method รับ `actor domain.User` เข้ามาเป็นพารามิเตอร์ตัวแรก
- `Place` บันทึกเจ้าของลงในออเดอร์ ส่วน `Get`, `Cancel`, `Total` ถามกฎในชั้น `domain` ก่อนทำงาน
- `SetStatus` ตรวจสิทธิ์ก่อนแตะฐานข้อมูล

```go
func (u OrderUseCase) Place(actor domain.User, menuID string, qty int) (domain.Order, error)
func (u OrderUseCase) Get(actor domain.User, id string) (domain.Order, error)
func (u OrderUseCase) Cancel(actor domain.User, id string) error
func (u OrderUseCase) SetStatus(actor domain.User, id string, next domain.Status) error
func (u OrderUseCase) Total(actor domain.User, id string) (float64, error)
```

ส่วนที่สำคัญที่สุดอยู่ใน `Cancel` ซึ่งแยกสองกรณีออกจากกันชัดเจน

```go
// ซ่อนออเดอร์ของผู้อื่นไว้หลัง 404 ด้วยเหตุผลเดียวกับ Get
if !order.CanBeViewedBy(actor) {
	return domain.ErrNotFound
}

// แยกสองกรณีให้ชัด เพราะความหมายต่างกันและ Status Code ก็ต่างกัน
//   ยกเลิกไม่ได้เพราะเริ่มทำอาหารแล้ว  -> 409
//   ยกเลิกไม่ได้เพราะไม่มีสิทธิ์        -> 403
if !order.CanCancel() {
	return domain.ErrNotCancelable
}
if !order.CanBeCancelledBy(actor) {
	return domain.ErrForbidden
}
```

**ไฟล์ที่ 2 — `usecase/order_test.go` สิ่งที่แก้ไข**

- ประกาศผู้ใช้สมมติสามคนไว้บนสุด แล้วเติม `actor` ให้ทุกจุดที่เรียก Use Case
- เพิ่ม Test ของกฎเรื่องสิทธิ์อีก 5 ตัว

```go
var (
	alice = domain.User{ID: "u1", Email: "alice@example.com", Role: domain.Customer}
	bob   = domain.User{ID: "u2", Email: "bob@example.com", Role: domain.Customer}
	staff = domain.User{ID: "u9", Email: "staff@example.com", Role: domain.Staff}
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
    // รับ actor เข้ามาเป็นพารามิเตอร์ตัวแรก เพื่อบันทึกว่าออเดอร์ใบนี้เป็นของผู้ใด
    func (u OrderUseCase) Place(actor domain.User, menuID string, qty int) (domain.Order, error) {
    	if qty < 1 {
    		return domain.Order{}, domain.ErrInvalidQty
    	}
    
    	item, err := u.repo.FindMenuItem(menuID)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("place order: %w", err)
    	}
    
    	order := domain.Order{
    		UserID: actor.ID,
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
    
    func (u OrderUseCase) Get(actor domain.User, id string) (domain.Order, error) {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return domain.Order{}, fmt.Errorf("get order %s: %w", id, err)
    	}
    
    	// คืน ErrNotFound ไม่ใช่ ErrForbidden โดยตั้งใจ
    	// เพราะการตอบ 403 เท่ากับยืนยันว่า "ออเดอร์ใบนี้มีอยู่จริง แต่ไม่ใช่ของคุณ"
    	// ซึ่งเป็นการเปิดเผยข้อมูลให้ผู้ที่ไล่เดาหมายเลขออเดอร์
    	if !order.CanBeViewedBy(actor) {
    		return domain.Order{}, domain.ErrNotFound
    	}
    	return order, nil
    }
    
    // ตัวอย่างที่ชัดเจนที่สุดว่าเหตุใดจึงต้องมีชั้น Use Case
    // เพราะไม่มี Constraint ใดในฐานข้อมูลที่เขียนกฎข้อนี้ได้
    func (u OrderUseCase) Cancel(actor domain.User, id string) error {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return fmt.Errorf("cancel order %s: %w", id, err)
    	}
    
    	// ซ่อนออเดอร์ของผู้อื่นไว้หลัง 404 ด้วยเหตุผลเดียวกับ Get
    	if !order.CanBeViewedBy(actor) {
    		return domain.ErrNotFound
    	}
    
    	// แยกสองกรณีให้ชัด เพราะความหมายต่างกันและ Status Code ก็ต่างกัน
    	//   ยกเลิกไม่ได้เพราะเริ่มทำอาหารแล้ว  -> 409
    	//   ยกเลิกไม่ได้เพราะไม่มีสิทธิ์        -> 403
    	if !order.CanCancel() {
    		return domain.ErrNotCancelable
    	}
    	if !order.CanBeCancelledBy(actor) {
    		return domain.ErrForbidden
    	}
    
    	order.Status = domain.Cancelled
    	return u.repo.Save(order)
    }
    
    func (u OrderUseCase) SetStatus(actor domain.User, id string, next domain.Status) error {
    	// การเดินสถานะเป็นงานของร้าน ตรวจสิทธิ์ก่อนแตะฐานข้อมูล
    	if !(domain.Order{}).CanChangeStatusBy(actor) {
    		return domain.ErrForbidden
    	}
    
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
    
    func (u OrderUseCase) Total(actor domain.User, id string) (float64, error) {
    	order, err := u.repo.FindOrder(id)
    	if err != nil {
    		return 0, fmt.Errorf("total %s: %w", id, err)
    	}
    	if !order.CanBeViewedBy(actor) {
    		return 0, domain.ErrNotFound
    	}
    	return order.Total() + domain.DeliveryFee, nil
    }
    ```
    

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `usecase/order_test.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    //
    // สังเกตว่ากฎเรื่องสิทธิ์ก็ทดสอบได้ที่ชั้นนี้เช่นกัน โดยไม่ต้องมี Token
    // และไม่ต้องส่ง HTTP Request เพราะสิทธิ์เป็นกฎทางธุรกิจ ไม่ใช่เรื่องของ HTTP
    
    var (
    	alice = domain.User{ID: "u1", Email: "alice@example.com", Role: domain.Customer}
    	bob   = domain.User{ID: "u2", Email: "bob@example.com", Role: domain.Customer}
    	staff = domain.User{ID: "u9", Email: "staff@example.com", Role: domain.Staff}
    )
    
    func TestCancelRefusedWhileCooking(t *testing.T) {
    	repo := &stubRepo{order: domain.Order{ID: "o1", UserID: alice.ID, Status: domain.Cooking}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	err := uc.Cancel(alice, "o1")
    
    	if !errors.Is(err, domain.ErrNotCancelable) {
    		t.Fatalf("ต้องการ ErrNotCancelable แต่ได้ %v", err)
    	}
    }
    
    func TestCancelAllowedWhilePending(t *testing.T) {
    	repo := &stubRepo{order: domain.Order{ID: "o1", UserID: alice.ID, Status: domain.Pending}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	if err := uc.Cancel(alice, "o1"); err != nil {
    		t.Fatalf("ไม่ควรเกิด error แต่ได้ %v", err)
    	}
    	if repo.saved.Status != domain.Cancelled {
    		t.Errorf("สถานะที่บันทึกควรเป็น cancelled แต่ได้ %v", repo.saved.Status)
    	}
    }
    
    func TestPlaceBelowMinimum(t *testing.T) {
    	repo := &stubRepo{item: domain.MenuItem{ID: "m1", Name: "Water", Price: 20}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	_, err := uc.Place(alice, "m1", 2) // 40 บาท ไม่ถึงยอดขั้นต่ำ 100 บาท
    
    	if !errors.Is(err, domain.ErrBelowMinimum) {
    		t.Fatalf("ต้องการ ErrBelowMinimum แต่ได้ %v", err)
    	}
    }
    
    func TestPlaceRejectsZeroQty(t *testing.T) {
    	repo := &stubRepo{item: domain.MenuItem{ID: "m1", Price: 500}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	if _, err := uc.Place(alice, "m1", 0); !errors.Is(err, domain.ErrInvalidQty) {
    		t.Fatalf("ต้องการ ErrInvalidQty แต่ได้ %v", err)
    	}
    }
    
    func TestPlaceWrapsNotFound(t *testing.T) {
    	uc := usecase.NewOrderUseCase(failingRepo{})
    
    	_, err := uc.Place(alice, "m9", 2)
    
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
    		UserID: alice.ID,
    		Status: domain.Pending,
    		Lines: []domain.OrderLine{
    			{Item: domain.MenuItem{Price: 80}, Qty: 2},
    			{Item: domain.MenuItem{Price: 120}, Qty: 1},
    		},
    	}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	got, err := uc.Total(alice, "o1")
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
    
    	if _, err := uc.Place(alice, "m1", 1); !errors.Is(err, domain.ErrBelowMinimum) {
    		t.Fatalf("ต้องการ ErrBelowMinimum แต่ได้ %v", err)
    	}
    
    	repo.item = domain.MenuItem{ID: "m2", Name: "Pad Thai", Price: 80}
    	order, err := uc.Place(alice, "m2", 2)
    	if err != nil {
    		t.Fatal(err)
    	}
    
    	if order.ID != "o1" {
    		t.Errorf("ออเดอร์แรกที่สำเร็จควรได้หมายเลข o1 แต่ได้ %q", order.ID)
    	}
    }
    
    // ─────────────── กฎเรื่องสิทธิ์ ทดสอบได้โดยไม่ต้องมี Token ───────────────
    
    func TestCustomerCannotCancelSomeoneElsesOrder(t *testing.T) {
    	repo := &stubRepo{order: domain.Order{ID: "o1", UserID: alice.ID, Status: domain.Pending}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	// bob เห็นออเดอร์ของ alice ไม่ได้ จึงต้องได้ ErrNotFound ไม่ใช่ ErrForbidden
    	if err := uc.Cancel(bob, "o1"); !errors.Is(err, domain.ErrNotFound) {
    		t.Fatalf("ต้องการ ErrNotFound แต่ได้ %v", err)
    	}
    	if repo.saved.Status == domain.Cancelled {
    		t.Error("ออเดอร์ของผู้อื่นต้องไม่ถูกยกเลิก")
    	}
    }
    
    func TestStaffCanCancelAnyOrder(t *testing.T) {
    	repo := &stubRepo{order: domain.Order{ID: "o1", UserID: alice.ID, Status: domain.Pending}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	if err := uc.Cancel(staff, "o1"); err != nil {
    		t.Fatalf("พนักงานควรยกเลิกได้ แต่ได้ %v", err)
    	}
    	if repo.saved.Status != domain.Cancelled {
    		t.Errorf("สถานะที่บันทึกควรเป็น cancelled แต่ได้ %v", repo.saved.Status)
    	}
    }
    
    func TestCustomerCannotChangeStatus(t *testing.T) {
    	repo := &stubRepo{order: domain.Order{ID: "o1", UserID: alice.ID, Status: domain.Pending}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	if err := uc.SetStatus(alice, "o1", domain.Confirmed); !errors.Is(err, domain.ErrForbidden) {
    		t.Fatalf("ต้องการ ErrForbidden แต่ได้ %v", err)
    	}
    }
    
    func TestGetHidesOtherPeoplesOrder(t *testing.T) {
    	repo := &stubRepo{order: domain.Order{ID: "o1", UserID: alice.ID, Status: domain.Pending}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	if _, err := uc.Get(bob, "o1"); !errors.Is(err, domain.ErrNotFound) {
    		t.Fatalf("ต้องการ ErrNotFound แต่ได้ %v", err)
    	}
    }
    
    func TestPlaceRecordsOwner(t *testing.T) {
    	repo := &stubRepo{item: domain.MenuItem{ID: "m1", Name: "Pad Thai", Price: 80}}
    	uc := usecase.NewOrderUseCase(repo)
    
    	order, err := uc.Place(alice, "m1", 2)
    	if err != nil {
    		t.Fatal(err)
    	}
    	if order.UserID != alice.ID {
    		t.Errorf("ออเดอร์ควรเป็นของ %s แต่ได้ %q", alice.ID, order.UserID)
    	}
    }
    ```
    

**ผลลัพธ์**

```bash
go test ./domain/ ./usecase/
```

```bash
ok  	food-api/domain	(cached)
ok  	food-api/usecase	0.412s
```

<aside>
🚨

ขั้นนี้ให้รัน `go test ./domain/ ./usecase/` **ไม่ใช่** `go build ./...`

เพราะการเปลี่ยนลายเซ็นของ Use Case ทำให้ `delivery/http/order.go` ที่ยังเรียกแบบเดิมอยู่ คอมไพล์ไม่ผ่านชั่วคราว

`not enough arguments in call to uc.Cancel`

**เป็นเรื่องปกติและจะหายไปในขั้นถัดไป** คอมไพเลอร์กำลังไล่จุดที่ต้องแก้ให้เราครบ เหมือนตอนขยาย `OrderRepository` ใน Day 2 ขั้นที่ 15

โปรเจกต์จะกลับมา `go build ./...` ผ่านทั้งหมดอีกครั้งที่ **ขั้นที่ 9**

</aside>

- `usecase/fakes_test.go` ไม่ต้องแก้เลย เพราะ `OrderRepository` ไม่ได้เปลี่ยน สิ่งที่เปลี่ยนคือลายเซ็นของ Use Case
- ถึงตรงนี้กฎเรื่องสิทธิ์ถูกทดสอบครบแล้ว ทั้งที่ยังไม่มีโค้ดเกี่ยวกับ Token อยู่ในโปรเจกต์แม้แต่บรรทัดเดียว

## 3.4 ชั้นที่ 3 — Delivery รับตัวตนเข้ามา

ชั้นนอกสุด เป็นชั้นเดียวที่รู้จักคำว่า Bearer

#### ขั้นที่ 5 — ชั้น Delivery — รับตัวตนเข้ามาทั้งชั้น

<aside>
✏️

**ขั้นนี้ทำชั้น Delivery ให้จบทั้งชั้น มีสองไฟล์ใหม่และสองไฟล์เดิม**

**1. สร้างไฟล์** `delivery/http/middleware_auth.go` แล้ววางโค้ดด้านล่างลงไป

**2. เปิดไฟล์** `delivery/http/order.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 10 แล้วแก้ตามรายการด้านล่าง

**3. รัน** `go build ./delivery/...` **ครั้งที่หนึ่ง** ตรวจสองไฟล์แรกก่อน

**4. เปิดไฟล์** `delivery/http/error.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 16 แล้วเพิ่ม `case` ตามด้านล่าง

**5. สร้างไฟล์** `delivery/http/auth.go` แล้ววางโค้ดด้านล่างลงไป

**6. รัน** `go build ./delivery/...` **ครั้งที่สอง**

ขั้นนี้มี **4 ไฟล์** แต่แบ่งเป็นสองช่วง ช่วงละสองไฟล์ และมีจุดตรวจของตนเอง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `delivery/http/middleware_auth.go` สร้างใหม่**

```go
package http

import (
	"strings"

	"food-api/domain"

	"github.com/gofiber/fiber/v2"
)

// actorKey คือชื่อช่องที่ใช้ฝากตัวตนของผู้เรียกไว้ใน Context ของ Request
const actorKey = "actor"

// RequireAuth เป็น Middleware ที่ทำหน้าที่เดียว คือแปลง Token เป็นตัวตน
//
// สังเกตว่าโค้ดที่รู้จักคำว่า Bearer มีอยู่ที่นี่ที่เดียวในโปรเจกต์
// ชั้น Use Case และชั้น Domain ไม่เคยเห็นคำนี้เลย
func (h Handler) RequireAuth(c *fiber.Ctx) error {
	raw := c.Get("Authorization")

	const prefix = "Bearer "
	if !strings.HasPrefix(raw, prefix) {
		return toHTTP(c, domain.ErrUnauthorized)
	}

	token := strings.TrimSpace(strings.TrimPrefix(raw, prefix))
	if token == "" {
		return toHTTP(c, domain.ErrUnauthorized)
	}

	actor, err := h.auth.Identify(token)
	if err != nil {
		return toHTTP(c, err)
	}

	// ฝากตัวตนไว้ให้ Handler ตัวถัดไปหยิบใช้
	c.Locals(actorKey, actor)
	return c.Next()
}

// RequireStaff ใช้ต่อจาก RequireAuth สำหรับเส้นทางที่เปิดให้เฉพาะพนักงาน
func (h Handler) RequireStaff(c *fiber.Ctx) error {
	if !actorOf(c).IsStaff() {
		return toHTTP(c, domain.ErrForbidden)
	}
	return c.Next()
}

// actorOf หยิบตัวตนออกจาก Context
// หากเส้นทางใดไม่ได้ผ่าน RequireAuth มาก่อน จะได้ User ว่างซึ่งไม่มีสิทธิ์ใด
func actorOf(c *fiber.Ctx) domain.User {
	u, _ := c.Locals(actorKey).(domain.User)
	return u
}
```

**ไฟล์ที่ 2 — `delivery/http/order.go` สิ่งที่แก้ไข**

- `Handler` ถือ `AuthUseCase` เพิ่มอีกตัว
- ทุก Handler ส่ง `actorOf(c)` เข้าไปให้ Use Case

```go
type Handler struct {
	uc   usecase.OrderUseCase
	auth usecase.AuthUseCase
}

func NewHandler(uc usecase.OrderUseCase, auth usecase.AuthUseCase) Handler {
	return Handler{uc: uc, auth: auth}
}
```

```go
order, err := h.uc.Place(actorOf(c), req.MenuID, req.Qty)
order, err := h.uc.Get(actorOf(c), c.Params("id"))
err := h.uc.Cancel(actorOf(c), c.Params("id"))
err := h.uc.SetStatus(actorOf(c), c.Params("id"), domain.Status(req.Status))
total, err := h.uc.Total(actorOf(c), id)
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `delivery/http/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/domain"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // Handler ถือ Use Case ไว้ ไม่ได้ถือการเชื่อมต่อฐานข้อมูล
    // บทนี้เพิ่ม AuthUseCase เข้ามาอีกตัว สำหรับงานพิสูจน์ตัวตน
    type Handler struct {
    	uc   usecase.OrderUseCase
    	auth usecase.AuthUseCase
    }
    
    func NewHandler(uc usecase.OrderUseCase, auth usecase.AuthUseCase) Handler {
    	return Handler{uc: uc, auth: auth}
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
    
    	order, err := h.uc.Place(actorOf(c), req.MenuID, req.Qty)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	return c.Status(201).JSON(toResponse(order))
    }
    
    func (h Handler) GetOrder(c *fiber.Ctx) error {
    	order, err := h.uc.Get(actorOf(c), c.Params("id"))
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(toResponse(order))
    }
    
    func (h Handler) CancelOrder(c *fiber.Ctx) error {
    	if err := h.uc.Cancel(actorOf(c), c.Params("id")); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    
    func (h Handler) SetStatus(c *fiber.Ctx) error {
    	var req setStatusRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	if err := h.uc.SetStatus(actorOf(c), c.Params("id"), domain.Status(req.Status)); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    
    func (h Handler) OrderTotal(c *fiber.Ctx) error {
    	id := c.Params("id")
    
    	total, err := h.uc.Total(actorOf(c), id)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	// ไม่มีการคำนวณใด ๆ ในชั้นนี้ การคำนวณทั้งหมดอยู่ใน domain และ usecase
    	return c.JSON(fiber.Map{"id": id, "total": total})
    }
    ```
    
- `middleware_auth.go` เป็นที่เดียวในโปรเจกต์ที่ปรากฏคำว่า `Bearer` และ `Authorization`
- `actorOf` คืน `User` ว่างเมื่อเส้นทางนั้นไม่ได้ผ่าน `RequireAuth` มาก่อน ซึ่ง `User` ว่างไม่ใช่ `staff` และไม่ได้เป็นเจ้าของออเดอร์ใด จึงไม่มีสิทธิ์ใดโดยปริยาย
- Handler ไม่ตัดสินใจเรื่องสิทธิ์เลย เพียงส่งต่อว่าใครเป็นผู้เรียก ค้นคำว่า `IsStaff` หรือ `OwnedBy` ในโฟลเดอร์ `delivery` จะไม่พบเลย

<aside>
💡

นี่คือ Decorator Pattern ตัวเดียวกับ Logger ที่เขียนไว้ใน Day 2 ขั้นที่ 19

งานที่ต้องทำเหมือนกันทุก Request เขียนไว้ที่เดียว แล้วครอบเส้นทางที่ต้องการ ไม่ใช่เขียนซ้ำในทุก Handler

</aside>

<aside>
💡

**ขั้นนี้มีสี่ไฟล์ แต่ให้รัน `go build ./delivery/...` สองครั้ง อย่ารวบไปรันทีเดียวตอนจบ**

สองไฟล์แรกต้องผ่านก่อน จึงค่อยทำสองไฟล์หลัง เพราะถ้าวางพลาดแล้วรันทีเดียวตอนจบ

ข้อความผิดพลาดจะชี้ไปที่ไฟล์หนึ่ง ทั้งที่ต้นเหตุอาจอยู่ในอีกไฟล์หนึ่ง และต้องไล่หาทั้งสี่ไฟล์

การตรวจครึ่งทางทำให้ขอบเขตของสิ่งที่ผิดแคบลงเหลือสองไฟล์เสมอ

</aside>

**ไฟล์ที่ 3 — `delivery/http/error.go` สิ่งที่แก้ไข**

```go
// เพิ่มในบทนี้ ยังคงเป็นไฟล์เดียวของทั้งโปรแกรมที่ตัดสินเรื่อง Status Code
case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidCredentials):
	// 401 แปลว่า "ยังไม่รู้ว่าคุณเป็นใคร"
	c.Set("WWW-Authenticate", "Bearer")
	return c.Status(401).JSON(msg(err.Error()))

case errors.Is(err, domain.ErrForbidden):
	// 403 แปลว่า "รู้แล้วว่าคุณเป็นใคร แต่คุณไม่มีสิทธิ์"
	return c.Status(403).JSON(msg(err.Error()))

case errors.Is(err, domain.ErrEmailTaken):
	return c.Status(409).JSON(msg(err.Error()))

case errors.Is(err, domain.ErrWeakPassword):
	return c.Status(422).JSON(msg(err.Error()))

case errors.Is(err, domain.ErrTooManyRequests):
	return c.Status(429).JSON(msg(err.Error()))
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `delivery/http/error.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    
    	// เพิ่มในบทนี้ ยังคงเป็นไฟล์เดียวของทั้งโปรแกรมที่ตัดสินเรื่อง Status Code
    	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidCredentials):
    		// 401 แปลว่า "ยังไม่รู้ว่าคุณเป็นใคร"
    		c.Set("WWW-Authenticate", "Bearer")
    		return c.Status(401).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrForbidden):
    		// 403 แปลว่า "รู้แล้วว่าคุณเป็นใคร แต่คุณไม่มีสิทธิ์"
    		return c.Status(403).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrEmailTaken):
    		return c.Status(409).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrWeakPassword):
    		return c.Status(422).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrTooManyRequests):
    		return c.Status(429).JSON(msg(err.Error()))
    
    	default:
    		// Error ที่ไม่รู้จัก หมายความว่าเป็นข้อผิดพลาดของระบบเราเอง
    		// บันทึกไว้ตรวจสอบภายหลัง แต่ไม่ส่งให้ผู้ใช้ เพราะอาจมีชื่อตารางปะปนอยู่
    		log.Println(err)
    		return c.Status(500).JSON(msg("internal error"))
    	}
    }
    ```
    

<aside>
💡

ไฟล์นี้โตขึ้นเป็นรอบที่สองแล้ว แต่ยังคงเป็น **จุดเดียวของทั้งโปรแกรม** ที่ตัดสินว่า Error ตัวใดควรเป็น Status Code ใด

ทั้งชั้น `domain` และชั้น `usecase` ยังไม่รู้จักเลขเหล่านี้เลยแม้แต่ตัวเดียว

</aside>

- `401` แถม Header `WWW-Authenticate: Bearer` ไปด้วย เป็นสิ่งที่มาตรฐาน HTTP กำหนดไว้ และช่วยให้เครื่องมืออย่าง Postman รู้ว่าต้องส่ง Token มาแบบใด

**ไฟล์ที่ 4 — `delivery/http/auth.go` สร้างใหม่**

```go
package http

import (
	"food-api/domain"

	"github.com/gofiber/fiber/v2"
)

// DTO ของการสมัครและเข้าสู่ระบบ อยู่ชั้นนอกเช่นเดิม
// registerRequest ตั้งใจไม่มีฟีลด์ Role
//
// ถ้ามีฟีลด์นี้ ผู้ใช้จะส่ง {"role":"staff"} เข้ามาแล้วตั้งตนเป็นพนักงานได้ทันที
// ช่องโหว่แบบนี้เรียกว่า Mass Assignment คือปล่อยให้ข้อมูลจากภายนอก
// ไหลเข้าไปตั้งค่าฟีลด์ที่ผู้ใช้ไม่ควรมีสิทธิ์กำหนด
type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userResponse ตั้งใจไม่มีฟีลด์ PasswordHash
// การแปลงผ่าน DTO ทำให้ค่าที่ไม่ควรออกไป ไม่มีทางหลุดออกไปโดยบังเอิญ
type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func toUserResponse(u domain.User) userResponse {
	return userResponse{ID: u.ID, Email: u.Email, Role: string(u.Role)}
}

func (h Handler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(msg("invalid body"))
	}

	u, err := h.auth.Register(req.Email, req.Password)
	if err != nil {
		return toHTTP(c, err)
	}
	return c.Status(201).JSON(toUserResponse(u))
}

func (h Handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(msg("invalid body"))
	}

	token, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		return toHTTP(c, err)
	}
	return c.JSON(tokenResponse{Token: token})
}

// Me ใช้ตรวจว่า Token ที่ถืออยู่แทนตัวตนใด
func (h Handler) Me(c *fiber.Ctx) error {
	return c.JSON(toUserResponse(actorOf(c)))
}
```

<aside>
🚨

สังเกต `userResponse` ให้ดี ไม่มีฟีลด์ `PasswordHash` อยู่เลย

หากส่ง `domain.User` ออกไปตรง ๆ ค่า Hash ของรหัสผ่านจะหลุดออกไปทาง API ทันที

นี่คือเหตุผลข้อที่สองของการมี DTO นอกเหนือจากการแยกรูปแบบข้อมูล คือ **กันไม่ให้สิ่งที่ไม่ควรออก หลุดออกไปโดยบังเอิญ**

</aside>

- Handler ทั้งสามตัวยังทำเพียงสามอย่างเช่นเดิม คือแปลง Request เรียก Use Case หนึ่งตัว แล้วเขียน Response

## 3.5 ค่าลับ และ Adapter ที่ทำงานจริง

ถึงตรงนี้ `usecase` ประกาศไว้แล้วว่าต้องการอะไร เหลือเพียงเขียนตัวที่ทำงานจริงมาเสียบ

#### ขั้นที่ 6 — ค่าลับของระบบยืนยันตัวตน ไว้ใน .env

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `.env.example` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 28

**2. เพิ่มสามค่าด้านล่างต่อท้ายไฟล์**

**3. คัดลอกลง `.env` ของตนเอง** `cp .env.example .env`

**4. สร้างค่าสุ่มของตนเอง** แล้วนำไปใส่เป็น `JWT_SECRET` ดูคำสั่งตามระบบปฏิบัติการด้านล่าง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

เพิ่มบรรทัดด้านล่างนี้ต่อท้ายไฟล์

```bash
# ─── เพิ่มในบทที่ 3 ค่าลับของระบบยืนยันตัวตน ───
# JWT_SECRET คือกุญแจที่ใช้เซ็นและตรวจ Token ผู้ใดได้ค่านี้ไป จะปลอม Token เป็นใครก็ได้
# จึงห้าม Commit ลง Git เด็ดขาด และต้องใช้คนละค่ากันระหว่างเครื่องพัฒนากับ Server จริง
#
# สร้างค่าสุ่มด้วยคำสั่ง   openssl rand -base64 48
JWT_SECRET=dev-only-secret-change-me-before-production-0123456789

# อายุของ Token สั้นไว้ปลอดภัยกว่า เพราะ Token ที่หลุดออกไปจะหมดอายุเร็ว
JWT_TTL=1h

# ต้นทางที่อนุญาตให้เรียก API นี้จากเบราว์เซอร์ ระบุให้ชัด ห้ามใช้ *
CORS_ORIGINS=http://localhost:5173
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `.env.example` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    
    # ─── เพิ่มในบทที่ 3 ค่าลับของระบบยืนยันตัวตน ───
    # JWT_SECRET คือกุญแจที่ใช้เซ็นและตรวจ Token ผู้ใดได้ค่านี้ไป จะปลอม Token เป็นใครก็ได้
    # จึงห้าม Commit ลง Git เด็ดขาด และต้องใช้คนละค่ากันระหว่างเครื่องพัฒนากับ Server จริง
    #
    # สร้างค่าสุ่ม   macOS, Linux, Git Bash:  openssl rand -base64 48
    # Windows PowerShell:  $b=[byte[]]::new(48); [Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($b); [Convert]::ToBase64String($b)
    JWT_SECRET=dev-only-secret-change-me-before-production-0123456789
    
    # อายุของ Token สั้นไว้ปลอดภัยกว่า เพราะ Token ที่หลุดออกไปจะหมดอายุเร็ว
    JWT_TTL=1h
    
    # ต้นทางที่อนุญาตให้เรียก API นี้จากเบราว์เซอร์ ระบุให้ชัด ห้ามใช้ *
    CORS_ORIGINS=http://localhost:5173
    
    # ─── วิธียืนยันตัวตน สลับได้ที่นี่ที่เดียว ───
    # jwt      เซิร์ฟเวอร์ไม่จำอะไรเลย ไม่ต้องอ่านฐานข้อมูล แต่ยกเลิก Token ก่อนหมดอายุไม่ได้
    # session  เซิร์ฟเวอร์จำไว้ในตาราง sessions ยกเลิกได้ทันที แต่ต้องอ่านฐานข้อมูลทุก Request
    AUTH_MODE=jwt
    
    # ─── บัญชีพนักงานคนแรก ───
    # เส้นทางสมัครสาธารณะสร้างได้เฉพาะลูกค้า บัญชีพนักงานจึงต้องมาจากผู้ดูแลระบบ
    # ถ้าเว้นว่างไว้ โปรแกรมจะไม่สร้างบัญชีพนักงานให้
    STAFF_EMAIL=staff@example.com
    STAFF_PASSWORD=password123
    ```
    
- `.env` ถูกกันไว้ด้วย `.gitignore` ตั้งแต่บทที่แล้ว จึงไม่ต้องแก้ `.gitignore` เพิ่มในบทนี้
- ค่าที่ให้มาในไฟล์ตัวอย่างเขียนไว้ชัดว่าเป็น `dev-only` เพื่อให้รันได้ทันทีในห้องเรียน แต่ต้องเปลี่ยนก่อนใช้งานจริงเสมอ

<aside>
🚨

`JWT_SECRET` มีค่าเท่ากับรหัสผ่านของผู้ใช้ทุกคนรวมกัน

ผู้ที่ได้กุญแจนี้ไป สร้าง Token ที่อ้างว่าเป็นผู้ใช้คนใดก็ได้ รวมทั้งอ้างว่าเป็น `staff` ด้วย

จึงห้ามเขียนค่านี้ลงในโค้ด ห้าม Commit ลง Git และต้องใช้คนละค่ากันระหว่างเครื่องพัฒนากับ Server จริง

</aside>

**สร้างค่าสุ่มสำหรับ `JWT_SECRET`**

macOS, Linux และ Git Bash บน Windows

```bash
openssl rand -base64 48
```

Windows PowerShell ซึ่งไม่มี `openssl` ติดมาให้

```powershell
$b = [byte[]]::new(48); [Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($b); [Convert]::ToBase64String($b)
```

<aside>
🚨

**ห้ามใช้ค่าตัวอย่างในไฟล์ และห้ามใช้ค่าเดียวกับเพื่อน**

ผู้ใดถือ `JWT_SECRET` ไป จะปลอม Token เป็นใครก็ได้ รวมถึงปลอมเป็น `staff`

ค่านี้ต้องยาวอย่างน้อย 32 ตัวอักษร โปรแกรมจะหยุดทำงานทันทีถ้าสั้นกว่านั้น

และต้องใช้คนละค่ากันระหว่างเครื่องพัฒนากับเครื่องจริงเสมอ

</aside>

<aside>
💡

ผู้ใช้ Windows ที่ยังไม่ได้อ่านหัวข้อ 3.0.4 แนะนำให้ย้อนไปอ่านก่อน เพราะยังมีคำสั่งอื่นในเอกสารนี้ที่ต้องเขียนต่างออกไป โดยเฉพาะ `curl` และการวนซ้ำในขั้นที่ 12

</aside>

#### ขั้นที่ 7 — Adapter ที่ทำงานจริง — bcrypt และ JWT

<aside>
📄

**1. สร้างโฟลเดอร์** `mkdir -p security`

**2. ติดตั้ง Library ทั้งสองตัวพร้อมกัน** `go get golang.org/x/crypto/bcrypt github.com/golang-jwt/jwt/v5`

**3. สร้างไฟล์ที่ 1** `security/bcrypt.go` แล้ววางโค้ดด้านล่างลงไป

**4. สร้างไฟล์ที่ 2** `security/jwt.go` แล้ววางโค้ดด้านล่างลงไป

**5. รัน** `go build ./security/`

ขั้นนี้มี **2 ไฟล์** คือ Adapter ทั้งสองตัวที่ชั้น Use Case ประกาศ Interface รอไว้แล้ว

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `security/bcrypt.go` สร้างใหม่**

```go
// Package security เก็บ Adapter ที่เชื่อมไลบรารีภายนอกเข้ากับ Interface ของชั้น Use Case
//
// อยู่วงเดียวกับ repository/ คือ Interface Adapters
// ต่างกันเพียงว่าเชื่อมกับไลบรารีเข้ารหัสและ Token ไม่ใช่ฐานข้อมูล
package security

import "golang.org/x/crypto/bcrypt"

// BcryptHasher ทำตาม Interface usecase.PasswordHasher
// โดยที่ชั้น Use Case ไม่รู้จักชื่อ bcrypt เลย
type BcryptHasher struct {
	cost int
}

// cost คือจำนวนรอบของการคำนวณ ยิ่งสูงยิ่งช้าและยิ่งเดายาก
// ค่าเริ่มต้นของไลบรารีคือ 10 ซึ่งเหมาะกับงานทั่วไป
func NewBcryptHasher(cost int) BcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return BcryptHasher{cost: cost}
}

func (h BcryptHasher) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	return string(b), err
}

// Compare คืน error เมื่อไม่ตรง ไม่ได้คืน bool
// เพราะ bcrypt เปรียบเทียบแบบใช้เวลาคงที่ ป้องกันการวัดเวลาเพื่อเดารหัสผ่าน
func (h BcryptHasher) Compare(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
```

<aside>
🧠

**เจาะลึก เหตุใดจึงไม่ใช้ SHA256**

SHA256 ถูกออกแบบมาให้ **เร็ว** ซึ่งเป็นคุณสมบัติที่ดีสำหรับการตรวจความถูกต้องของไฟล์ แต่เป็นคุณสมบัติที่แย่มากสำหรับรหัสผ่าน

เพราะการ์ดจอหนึ่งใบคำนวณ SHA256 ได้หลายพันล้านครั้งต่อวินาที ผู้โจมตีที่ได้ฐานข้อมูลไปจึงไล่เดารหัสผ่านได้เร็วมาก

`bcrypt` ถูกออกแบบให้ **ช้าโดยตั้งใจ** และปรับความช้าได้ด้วยค่า Cost ยิ่งค่าสูงยิ่งใช้เวลาคำนวณนานขึ้นแบบทวีคูณ

นอกจากนี้ `bcrypt` ยังใส่ค่าสุ่มที่เรียกว่า Salt ลงไปให้เองในทุกครั้ง ผู้ใช้สองคนที่ตั้งรหัสผ่านเหมือนกัน จึงได้ค่า Hash ที่ต่างกัน

</aside>

- `Compare` คืน `error` ไม่ได้คืน `bool` เพราะไลบรารีเปรียบเทียบแบบใช้เวลาคงที่ ป้องกันการวัดเวลาตอบกลับเพื่อค่อย ๆ เดารหัสผ่านทีละตัวอักษร

**ไฟล์ที่ 2 — `security/jwt.go` สร้างใหม่**

```go
package security

import (
	"errors"
	"time"

	"food-api/domain"

	"github.com/golang-jwt/jwt/v5"
)

// JWTIssuer ทำตาม Interface usecase.TokenIssuer
//
// ไฟล์นี้เป็นไฟล์เดียวในโปรเจกต์ที่ import ไลบรารี JWT
// หากวันหนึ่งเปลี่ยนไปใช้ระบบ Token แบบอื่น จะแก้เฉพาะไฟล์นี้ไฟล์เดียว
type JWTIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTIssuer(secret string, ttl time.Duration) JWTIssuer {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return JWTIssuer{secret: []byte(secret), ttl: ttl}
}

func (j JWTIssuer) Issue(u domain.User) (string, error) {
	now := time.Now()

	// Payload ของ JWT อ่านได้ด้วยตาเปล่า ห้ามใส่รหัสผ่านหรือข้อมูลลับลงไป
	// ใส่เพียงพอให้รู้ว่าเป็นใครและมีสิทธิ์ระดับใด
	claims := jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"role":  string(u.Role),
		"iat":   now.Unix(),
		"exp":   now.Add(j.ttl).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
}

func (j JWTIssuer) Parse(token string) (domain.User, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// ต้องตรวจว่าอัลกอริทึมเป็นตัวที่เราตั้งใจใช้จริง
		// มิฉะนั้นผู้โจมตีอาจส่ง Token ที่ระบุ alg เป็น none เข้ามา
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil || !t.Valid {
		return domain.User{}, domain.ErrUnauthorized
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return domain.User{}, domain.ErrUnauthorized
	}

	id, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)
	if id == "" {
		return domain.User{}, domain.ErrUnauthorized
	}

	return domain.User{ID: id, Email: email, Role: domain.Role(role)}, nil
}
```

- ไฟล์นี้เป็น **ไฟล์เดียวในโปรเจกต์** ที่ import ไลบรารี JWT ค้นทั้งโปรเจกต์จะไม่พบคำนี้ที่อื่นอีกเลย
- `Issue` ใส่ลง Payload เพียงเท่าที่จำเป็น คือ `sub`, `email`, `role` และเวลาหมดอายุ ไม่ใส่รหัสผ่านหรือข้อมูลลับใด

<aside>
🚨

บรรทัดที่ตรวจ `SigningMethodHMAC` และ `jwt.WithValidMethods` สำคัญมาก ห้ามตัดออก

หากไม่ตรวจว่าอัลกอริทึมเป็นตัวที่เราตั้งใจใช้จริง ผู้โจมตีจะส่ง Token ที่ระบุ `alg` เป็น `none` เข้ามา แล้วผ่านการตรวจไปได้โดยไม่ต้องมีลายเซ็น

ช่องโหว่นี้มีชื่อเรียกว่า **alg=none attack** และเคยทำให้ไลบรารี JWT หลายตัวถูกเจาะมาแล้วในอดีต

</aside>

#### ขั้นที่ 8 — ที่เก็บผู้ใช้ ทั้งแบบ Memory และ Postgres

<aside>
📄

**1. สร้างไฟล์ที่ 1** `repository/memory/user.go` แล้ววางโค้ดด้านล่างลงไป

**2. สร้างไฟล์ที่ 2** `repository/postgres/user.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./repository/...`

ขั้นนี้มี **2 ไฟล์** สร้างให้ครบทุกไฟล์ก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

ไฟล์ `repository/memory/user.go`

```go
package memory

import (
	"strings"
	"sync"

	"food-api/domain"
)

// UserRepo ทำตาม Interface usecase.UserRepository
// รูปแบบเดียวกับ OrderRepo ในบทที่แล้ว คือเก็บและอ่านข้อมูลเท่านั้น ไม่มีกฎทางธุรกิจ
type UserRepo struct {
	mu    *sync.RWMutex
	users map[string]domain.User // key คืออีเมล
}

func NewUserRepo() UserRepo {
	return UserRepo{mu: &sync.RWMutex{}, users: map[string]domain.User{}}
}

func (r UserRepo) CreateUser(u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[strings.ToLower(u.Email)] = u
	return nil
}

func (r UserRepo) FindUserByEmail(email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[strings.ToLower(email)]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}
```

ไฟล์ `repository/postgres/user.go`

```go
package postgres

import (
	"context"
	"errors"

	"food-api/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepo {
	return UserRepo{db: db}
}

func (r UserRepo) CreateUser(u domain.User) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)`,
		u.ID, u.Email, u.PasswordHash, string(u.Role))

	// 23505 คือรหัสของ Postgres ที่แปลว่า "ค่าซ้ำกับที่มีอยู่แล้ว"
	// แปลงเป็น Error ของ Domain ที่ชั้นนี้ เพื่อไม่ให้ชั้นในต้องรู้จักรหัสของ Postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_pkey" {
		return domain.ErrIDTaken
	}
	return err
}

func (r UserRepo) FindUserByEmail(email string) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, email, password_hash, role FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)

	// แปลง Error ของ Driver เป็น Error ของ Domain ที่ชั้นนี้ เช่นเดิม
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	return u, err
}
```

- รูปแบบเหมือน `OrderRepo` ของบทที่แล้วทุกประการ คือเก็บและอ่านข้อมูลเท่านั้น ไม่มีกฎทางธุรกิจ
- `postgres.UserRepo` แปลง `pgx.ErrNoRows` เป็น `domain.ErrNotFound` ที่ชั้นนี้เช่นเดิม Error ของ Driver จึงไม่รั่วขึ้นไปยังชั้นบน

## 3.6 ประกอบเข้าด้วยกัน แล้วสมัครสมาชิกครั้งแรก

ชิ้นส่วนครบแล้ว เหลือเพียงต่อสายใน `main.go` ขั้นนี้จะกลับมา `go build ./...` ผ่านทั้งโปรเจกต์อีกครั้ง

#### ขั้นที่ 9 — ประกอบเข้าด้วยกัน แล้วสมัครสมาชิกครั้งแรก

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 11

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go build ./...` แล้วตามด้วย `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- อ่าน `JWT_SECRET` และ `JWT_TTL` จาก Environment แล้วหยุดทำงานทันทีเมื่อไม่มี
- สร้าง Adapter ทั้งสองตัว แล้วส่งเข้าไปให้ `AuthUseCase`
- แบ่ง Route ออกเป็นสามกลุ่ม คือเปิดให้ทุกคน ต้องเข้าสู่ระบบ และเฉพาะพนักงาน

```go
// ค่าลับต้องมาจาก Environment เสมอ และต้องหยุดทำงานทันทีเมื่อไม่มี
// ห้ามใส่ค่าสำรองไว้ในโค้ด เพราะจะกลายเป็นรหัสลับที่ทุกคนบนโลกรู้
secret := os.Getenv("JWT_SECRET")
if len(strings.TrimSpace(secret)) < 32 {
	log.Fatal("ต้องกำหนด JWT_SECRET ใน .env และต้องยาวอย่างน้อย 32 ตัวอักษร")
}
```

การแบ่ง Route เป็นสามกลุ่ม ทำให้อ่านไฟล์เดียวก็รู้ว่าเส้นทางใดต้องมีสิทธิ์ระดับใด

```go
// ── เปิดให้ทุกคน ──
api.Post("/auth/register", h.Register)
api.Post("/auth/login", h.Login)
api.Get("/menu", h.ListMenu)

// ── ต้องเข้าสู่ระบบก่อน ──
me := api.Group("", h.RequireAuth)
me.Get("/auth/me", h.Me)
me.Post("/orders", h.PlaceOrder)
me.Get("/orders/:id", h.GetOrder)
me.Get("/orders/:id/total", h.OrderTotal)
me.Post("/orders/:id/cancel", h.CancelOrder)

// ── เฉพาะพนักงานร้าน ──
staff := api.Group("", h.RequireAuth, h.RequireStaff)
staff.Patch("/orders/:id/status", h.SetStatus)
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"context"
    	"log"
    	"os"
    	"strings"
    	"time"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/repository/postgres"
    	"food-api/security"
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
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึงทุกชั้นพร้อมกัน และเป็นที่เดียวที่เลือกว่า
    // จะใช้ Adapter ตัวใด ทั้งที่เก็บข้อมูล วิธี Hash รหัสผ่าน และรูปแบบ Token
    func main() {
    	_ = godotenv.Load()
    
    	// ค่าลับต้องมาจาก Environment เสมอ และต้องหยุดทำงานทันทีเมื่อไม่มี
    	// ห้ามใส่ค่าสำรองไว้ในโค้ด เพราะจะกลายเป็นรหัสลับที่ทุกคนบนโลกรู้
    	secret := os.Getenv("JWT_SECRET")
    	if len(strings.TrimSpace(secret)) < 32 {
    		log.Fatal("ต้องกำหนด JWT_SECRET ใน .env และต้องยาวอย่างน้อย 32 ตัวอักษร")
    	}
    
    	ttl, err := time.ParseDuration(envOr("JWT_TTL", "1h"))
    	if err != nil {
    		log.Fatal("JWT_TTL ไม่ถูกต้อง: ", err)
    	}
    
    	// 1. ที่เก็บข้อมูล เลือกจาก Environment เช่นเดิม
    	var orders usecase.OrderRepository
    	var users usecase.UserRepository
    
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
    		orders = postgres.NewOrderRepo(pool)
    		users = postgres.NewUserRepo(pool)
    	} else {
    		log.Println("ใช้ Memory")
    		orders = memory.NewOrderRepo(seedMenu())
    		users = memory.NewUserRepo()
    	}
    
    	// 2. Adapter ด้านความปลอดภัย
    	//    ชั้น Use Case รู้จักเพียง Interface ไม่รู้จักชื่อ bcrypt และ JWT
    	hasher := security.NewBcryptHasher(0)
    	tokens := security.NewJWTIssuer(secret, ttl)
    
    	// 3. กฎ
    	uc := usecase.NewOrderUseCase(orders)
    	auth := usecase.NewAuthUseCase(users, hasher, tokens)
    
    	// บัญชีพนักงานคนแรกมาจากผู้ดูแลระบบ ไม่ใช่จากการสมัครเอง
    	// เส้นทาง /auth/register สาธารณะสร้างได้เฉพาะลูกค้าเท่านั้น
    	if email, pw := os.Getenv("STAFF_EMAIL"), os.Getenv("STAFF_PASSWORD"); email != "" && pw != "" {
    		if _, err := auth.CreateStaff(email, pw); err != nil {
    			log.Println("ข้ามการสร้างบัญชีพนักงานเริ่มต้น:", err)
    		} else {
    			log.Println("สร้างบัญชีพนักงานเริ่มต้นแล้ว:", email)
    		}
    	}
    
    	// 4. ชั้นขนส่ง
    	h := http.NewHandler(uc, auth)
    
    	app := fiber.New()
    
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	api := app.Group("/api")
    
    	// ── เปิดให้ทุกคน ──
    	api.Post("/auth/register", h.Register)
    	api.Post("/auth/login", h.Login)
    	api.Get("/menu", h.ListMenu)
    
    	// ── ต้องเข้าสู่ระบบก่อน ──
    	me := api.Group("", h.RequireAuth)
    	me.Get("/auth/me", h.Me)
    	me.Post("/orders", h.PlaceOrder)
    	me.Get("/orders/:id", h.GetOrder)
    	me.Get("/orders/:id/total", h.OrderTotal)
    	me.Post("/orders/:id/cancel", h.CancelOrder)
    
    	// ── เฉพาะพนักงานร้าน ──
    	staff := api.Group("", h.RequireAuth, h.RequireStaff)
    	staff.Patch("/orders/:id/status", h.SetStatus)
    
    	log.Fatal(app.Listen(":3000"))
    }
    
    func envOr(key, fallback string) string {
    	if v := os.Getenv(key); v != "" {
    		return v
    	}
    	return fallback
    }
    ```
    

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -i -X POST localhost:3000/api/auth/register -H 'Content-Type: application/json' -d '{"email":"alice@example.com","password":"password123"}'
```

```bash
HTTP/1.1 201 Created
Content-Type: application/json

{"id":"u1","email":"alice@example.com","role":"customer"}
```

คำสั่งที่ 2

```bash
curl -s -X POST localhost:3000/api/auth/login -H 'Content-Type: application/json' -d '{"email":"alice@example.com","password":"password123"}'
```

```bash
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaWNlQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY4MDM4NzUzLCJpYXQiOjE3NjgwMzUxNTMsInJvbGUiOiJjdXN0b21lciIsInN1YiI6InUxIn0.Yl0hSJ7v0_0Xr7zVQ0Yt6yqYqW9cVYb7yYh8mQ2xTnA"}
```

คำสั่งที่ 3

```bash
curl -i -X POST localhost:3000/api/orders -H 'Content-Type: application/json' -d '{"menu_id":"m1","qty":2}'
```

```bash
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Bearer
Content-Type: application/json

{"error":"unauthorized"}
```

<aside>
💡

คำสั่งที่ 3 คือช่องโหว่ที่เปิดบทมาถูกปิดแล้ว การสั่งงานออเดอร์โดยไม่แนบ Token ได้ `401` ทันที

ให้เก็บค่า `token` จากคำสั่งที่ 2 ไว้ในตัวแปรของ Shell เพื่อใช้ในขั้นถัดไป

</aside>

```bash
TOKEN=$(curl -s -X POST localhost:3000/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123"}' | jq -r .token)

echo $TOKEN
```

<aside>
🚨

หากโปรแกรมหยุดทันทีพร้อมข้อความ `ต้องกำหนด JWT_SECRET ใน .env` แสดงว่ายังไม่ได้ทำขั้นที่ 6 ให้ครบ หรือลืมคัดลอก `.env.example` เป็น `.env`

การหยุดทำงานเสียงดังแบบนี้เป็นความตั้งใจ ดีกว่าปล่อยให้ระบบรันต่อไปด้วยกุญแจที่คาดเดาได้

</aside>

## 3.7 ย้ายผู้ใช้ขึ้น Postgres

จนถึงขั้นที่แล้วยังเก็บผู้ใช้ไว้ใน Memory ซึ่งหายทุกครั้งที่ปิดโปรแกรม ขั้นตอนถัดไปคือย้ายขึ้นฐานข้อมูลจริง

#### ขั้นที่ 10 — ย้ายผู้ใช้ขึ้น Postgres — Migration ที่ 4 และ 5 พร้อม Adapter

<aside>
✏️

**ขั้นนี้ย้ายผู้ใช้ขึ้น Postgres ให้จบในขั้นเดียว มีสองไฟล์ใหม่และหนึ่งไฟล์เดิม**

**1. สร้างไฟล์ที่ 1** `migrations/004_users.sql` แล้ววางโค้ดด้านล่างลงไป

**2. สร้างไฟล์ที่ 2** `migrations/005_orders_user_id.sql` แล้ววางโค้ดด้านล่างลงไป

**3. เปิดไฟล์** `repository/postgres/order.go` ที่สร้างไว้ตั้งแต่ Day 2 ขั้นที่ 35 แล้วแก้ตามรายการด้านล่าง

**4. เปิดฐานข้อมูล** `docker compose up -d` แล้วรอจนขึ้น `(healthy)`

**5. รัน** `go run ./cmd/migrate` แล้วตามด้วย `go build ./...`

ขั้นนี้มี **3 ไฟล์** สร้างและแก้ให้ครบก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `migrations/004_users.sql` สร้างใหม่**

```sql
-- ตารางผู้ใช้ เก็บเฉพาะรหัสผ่านที่ผ่านการ Hash แล้ว
-- ไม่มีคอลัมน์ใดในระบบที่เก็บรหัสผ่านจริง

CREATE TABLE users (
    id            TEXT PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL DEFAULT 'customer'
                  CHECK (role IN ('customer', 'staff')),
    created       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ไม่ต้องสร้าง Index ให้ email เพิ่ม เพราะ UNIQUE สร้าง Index ให้อยู่แล้ว
```

**ไฟล์ที่ 2 — `migrations/005_orders_user_id.sql` สร้างใหม่**

```sql
-- ผูกออเดอร์เข้ากับเจ้าของ
--
-- ต้องเป็น Nullable เพราะออเดอร์ที่มีอยู่แล้วในฐานข้อมูลยังไม่มีเจ้าของ
-- เป็นกฎเดียวกับที่ใช้ตอนเพิ่มคอลัมน์ note ในบทที่แล้ว
ALTER TABLE orders ADD COLUMN user_id TEXT REFERENCES users(id);

-- หน้าจอ "ออเดอร์ของฉัน" ค้นด้วยคอลัมน์นี้บ่อย จึงต้องมี Index
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

<aside>
💡

ผลลัพธ์ขึ้นว่า **2 ไฟล์** ไม่ใช่ 5 ไฟล์ เพราะ `001` ถึง `003` ถูกบันทึกไว้ใน `schema_migrations` ตั้งแต่บทที่แล้ว

ข้อมูลเมนูและออเดอร์เดิมจึงยังอยู่ครบ ไม่ถูกแตะต้อง นี่คือเหตุผลทั้งหมดที่บทที่แล้วต้องเขียนตัวรัน Migration ขึ้นมา

</aside>

- `user_id` ใน `005` ประกาศเป็น **Nullable** เพราะออเดอร์ที่มีอยู่แล้วในฐานข้อมูลยังไม่มีเจ้าของ เป็นกฎเดียวกับตอนเพิ่มคอลัมน์ `note` ใน Day 2 ขั้นที่ 33
- ไม่ได้สร้าง Index ให้ `email` เพิ่ม เพราะ `UNIQUE` สร้าง Index ให้อยู่แล้ว สังเกตบรรทัด `users_email_key` ในผลลัพธ์ด้านบน

**ไฟล์ที่ 3 — `repository/postgres/order.go` สิ่งที่แก้ไข**

- `Save` เขียนคอลัมน์ `user_id` เพิ่ม
- `FindOrder` อ่านคอลัมน์ `user_id` กลับมา และต้องรองรับค่า `NULL` ของออเดอร์เก่า

```go
var o domain.Order
var userID *string // ออเดอร์เก่าก่อน migration 005 ยังไม่มีเจ้าของ จึงต้องรับ NULL ได้
err := r.db.QueryRow(ctx,
	`SELECT id, status, user_id FROM orders WHERE id = $1`, id,
).Scan(&o.ID, &o.Status, &userID)
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `repository/postgres/order.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    	var userID *string // ออเดอร์เก่าก่อน migration 005 ยังไม่มีเจ้าของ จึงต้องรับ NULL ได้
    	err := r.db.QueryRow(ctx,
    		`SELECT id, status, user_id FROM orders WHERE id = $1`, id,
    	).Scan(&o.ID, &o.Status, &userID)
    
    	if errors.Is(err, pgx.ErrNoRows) {
    		return domain.Order{}, domain.ErrNotFound
    	}
    	if err != nil {
    		return domain.Order{}, err
    	}
    	if userID != nil {
    		o.UserID = *userID
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
    		INSERT INTO orders (id, status, total, user_id) VALUES ($1, $2, $3, $4)
    		ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status,
    		                               total  = EXCLUDED.total`,
    		o.ID, string(o.Status), o.Total(), o.UserID)
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
    

<aside>
⚠️

ใช้ `*string` ไม่ใช่ `string` เพราะคอลัมน์เป็น Nullable

หากประกาศเป็น `string` ธรรมดา การอ่านออเดอร์เก่าที่ `user_id` เป็น `NULL` จะล้มเหลวทันทีด้วยข้อความ `cannot scan NULL into *string`

นี่คือราคาที่ต้องจ่ายของกฎ *เพิ่มคอลัมน์แบบ Nullable เสมอ* ซึ่งคุ้มกว่าการทำให้ระบบที่มีข้อมูลอยู่แล้วใช้งานไม่ได้

</aside>

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
go run ./cmd/migrate
```

```bash
ใช้งาน 004_users แล้ว
ใช้งาน 005_orders_user_id แล้ว
ใช้งาน migration ใหม่ทั้งหมด 2 ไฟล์
```

คำสั่งที่ 2

```bash
docker compose exec db psql -U app -d food -c "\d users"
```

```bash
 Column        |           Type           | Nullable |     Default
---------------+--------------------------+----------+------------------
 id            | text                     | not null |
 email         | text                     | not null |
 password_hash | text                     | not null |
 role          | text                     | not null | 'customer'::text
 created       | timestamp with time zone | not null | now()
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
```

### 3.7.1 ทำไมทุก Migration จึงมีบรรทัด CREATE INDEX

ใน `005_orders_user_id.sql` มีสองบรรทัดที่ทำงานคนละหน้าที่กัน บรรทัดแรกเพิ่มคอลัมน์ ส่วนบรรทัดที่สองสร้าง Index

```sql
ALTER TABLE orders ADD COLUMN IF NOT EXISTS user_id TEXT;
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders (user_id);
```

**Index คือสำเนาของคอลัมน์ที่เรียงลำดับไว้ล่วงหน้า พร้อมตัวชี้กลับไปยังแถวจริง**

เทียบกับหนังสือหนึ่งเล่ม การหาคำว่า Middleware มีสองวิธี วิธีแรกคือเปิดอ่านทีละหน้าตั้งแต่หน้าแรก วิธีที่สองคือเปิดดัชนีท้ายเล่มซึ่งเรียงตามตัวอักษรไว้แล้ว แล้วกระโดดไปหน้าที่ระบุ ฐานข้อมูลทำงานแบบเดียวกันทุกประการ

|  | ไม่มี Index | มี Index |
| --- | --- | --- |
| ชื่อแผนที่ Postgres ใช้ | `Seq Scan` | `Index Scan` |
| วิธีทำงาน | อ่านทุกแถวในตาราง แล้วคัดทิ้งทีละแถว | กระโดดตรงไปยังแถวที่ตรงเงื่อนไข |
| ตาราง 10 แถว | เร็วพอกัน ต่างกันไม่ถึงมิลลิวินาที | เร็วพอกัน |
| ตาราง 5 ล้านแถว | อ่านครบ 5 ล้านแถวทุกครั้งที่เรียก | อ่านไม่กี่แถว |

**เหตุใดจึงเป็นคอลัมน์ `user_id`**

เพราะตั้งแต่บทนี้เป็นต้นไป คำถามที่ระบบถามฐานข้อมูลบ่อยที่สุดคือ *ออเดอร์เหล่านี้เป็นของใคร* ทุกครั้งที่ตรวจสิทธิ์ และทุกครั้งที่ลูกค้าเปิดดูรายการออเดอร์ของตนเอง คำสั่งที่วิ่งไปคือ

```sql
SELECT ... FROM orders WHERE user_id = $1;
```

คอลัมน์ที่ปรากฏหลัง `WHERE` บ่อย ๆ คือคอลัมน์ที่คุ้มค่าแก่การสร้าง Index

<aside>
⚠️

**Index ไม่ได้ฟรี จึงห้ามสร้างพร่ำเพรื่อ**

ทุกครั้งที่ `INSERT`, `UPDATE` หรือ `DELETE` ฐานข้อมูลต้องแก้ทั้งตารางจริงและทุก Index ที่เกี่ยวข้อง การเขียนจึงช้าลงตามจำนวน Index

Index ยังกินพื้นที่ดิสก์เพิ่มอีกชุดหนึ่งด้วย

หลักที่ใช้ตัดสินคือ **สร้าง Index ให้คอลัมน์ที่ใช้ค้นหาบ่อย บนตารางที่อ่านมากกว่าเขียน** ไม่ใช่สร้างให้ทุกคอลัมน์

</aside>

**สามกรณีที่ไม่ต้องสร้างเพิ่ม**

- คอลัมน์ที่เป็น `PRIMARY KEY` มี Index ให้อยู่แล้วโดยอัตโนมัติ เช่น `orders.id`
- คอลัมน์ที่ประกาศ `UNIQUE` ก็มีให้แล้วเช่นกัน เช่น `users.email` สังเกตบรรทัด `users_email_key` ในผลลัพธ์ของขั้นที่ 10
- คอลัมน์ที่มีค่าซ้ำกันมากจนแทบไม่ช่วยกรอง เช่น คอลัมน์ที่มีเพียงค่า `true` กับ `false`

ส่วน `006_sessions.sql` ในหัวข้อ 3.11 สร้าง Index ให้ `expires_at` ด้วยเหตุผลเดียวกัน คือคำสั่งกวาดล้าง Session ที่หมดอายุจะถามด้วยคอลัมน์นี้

**เห็นด้วยตาว่าต่างกันจริงหรือไม่ ใช้ `EXPLAIN ANALYZE`**

คำสั่งนี้บอกให้ Postgres เล่าว่ามันวางแผนทำงานอย่างไร และใช้เวลาไปเท่าใด ตัวเลขด้านล่างมาจากตาราง `orders` ที่ใส่ข้อมูลไว้สองแสนแถว

```sql
EXPLAIN ANALYZE SELECT id, status, total FROM orders WHERE user_id = 'x42';
```

**เมื่อมี Index**

```bash
Bitmap Heap Scan on orders  (cost=4.60..147.55 rows=40 width=23) (actual time=0.070..0.749 rows=40 loops=1)
  Recheck Cond: (user_id = 'x42'::text)
  ->  Bitmap Index Scan on idx_orders_user_id  (actual time=0.053..0.053 rows=40 loops=1)
        Index Cond: (user_id = 'x42'::text)
Planning Time: 0.800 ms
Execution Time: 0.841 ms
```

**เมื่อลบ Index ทิ้ง**

```bash
Seq Scan on orders  (cost=0.00..4369.00 rows=40 width=23) (actual time=0.009..24.078 rows=40 loops=1)
  Filter: (user_id = 'x42'::text)
  Rows Removed by Filter: 199960
Planning Time: 0.888 ms
Execution Time: 24.129 ms
```

<aside>
💡

บรรทัดที่ต้องอ่านให้ออกคือ **`Rows Removed by Filter: 199960`**

แปลว่าฐานข้อมูลอ่านมาสองแสนแถว เพื่อจะเอาไปใช้จริงเพียง 40 แถว ที่เหลือคือแรงที่เสียเปล่าทั้งหมด

ผลคือ **24.1 มิลลิวินาที เทียบกับ 0.8 มิลลิวินาที** ต่างกันราวสามสิบเท่า และช่องว่างนี้จะถ่างขึ้นเรื่อย ๆ เมื่อข้อมูลโตขึ้น

เมื่อเห็นคำว่า `Seq Scan` บนตารางใหญ่ในเส้นทางที่เรียกบ่อย นั่นคือสัญญาณว่าควรพิจารณาสร้าง Index

</aside>

## 3.8 Rate Limiting และ Hardening

ระบบตอนนี้ปลอดภัยขึ้นแล้ว แต่ยังมีปัญหาหนึ่งข้อ คือผู้โจมตียิงหน้าเข้าสู่ระบบกี่ครั้งก็ได้ เพื่อไล่เดารหัสผ่าน

#### ขั้นที่ 11 — Middleware สองตัว — จำกัดจำนวน Request และ Security Header

<aside>
📄

**1. สร้างไฟล์ที่ 1** `delivery/http/middleware_ratelimit.go` แล้ววางโค้ดด้านล่างลงไป

**2. สร้างไฟล์ที่ 2** `delivery/http/middleware_headers.go` แล้ววางโค้ดด้านล่างลงไป

**3. รัน** `go build ./...`

ขั้นนี้มี **2 ไฟล์** ทั้งคู่เป็น Middleware ที่เขียนไว้ก่อน แต่ยังไม่ถูกใช้งานจนกว่าจะติดตั้งในขั้นถัดไป

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `delivery/http/middleware_ratelimit.go` สร้างใหม่**

```go
package http

import (
	"fmt"
	"sync"
	"time"

	"food-api/domain"

	"github.com/gofiber/fiber/v2"
)

// RateLimiter นับจำนวน Request ต่อหนึ่งช่วงเวลา แบบ Fixed Window
//
// เขียนเองประมาณ 40 บรรทัด เพื่อให้เห็นว่าเบื้องหลังคือ map กับตัวนับเท่านั้น
// และเพื่อให้เห็นชัดว่าตัวนับนี้อยู่ในหน่วยความจำของ Process
type RateLimiter struct {
	mu     *sync.Mutex
	hits   map[string]*window
	limit  int
	period time.Duration
}

type window struct {
	count   int
	resetAt time.Time
}

func NewRateLimiter(limit int, period time.Duration) *RateLimiter {
	return &RateLimiter{
		mu:     &sync.Mutex{},
		hits:   map[string]*window{},
		limit:  limit,
		period: period,
	}
}

func (rl *RateLimiter) Limit(c *fiber.Ctx) error {
	key := c.IP() // นับแยกตามหมายเลข IP ของผู้เรียก

	rl.mu.Lock()
	now := time.Now()
	w, ok := rl.hits[key]
	if !ok || now.After(w.resetAt) {
		w = &window{resetAt: now.Add(rl.period)}
		rl.hits[key] = w
	}
	w.count++
	count, resetAt := w.count, w.resetAt
	rl.mu.Unlock()

	remaining := rl.limit - count
	if remaining < 0 {
		remaining = 0
	}
	c.Set("X-RateLimit-Limit", fmt.Sprint(rl.limit))
	c.Set("X-RateLimit-Remaining", fmt.Sprint(remaining))

	if count > rl.limit {
		// Retry-After บอกผู้เรียกว่าให้รออีกกี่วินาที เป็นมารยาทที่ควรมีเสมอ
		c.Set("Retry-After", fmt.Sprint(int(time.Until(resetAt).Seconds())+1))
		return toHTTP(c, domain.ErrTooManyRequests)
	}
	return c.Next()
}
```

- วิธีนับแบบนี้เรียกว่า Fixed Window คือแบ่งเวลาเป็นช่วง ๆ แล้วนับจำนวนครั้งในแต่ละช่วง เมื่อหมดช่วงก็เริ่มนับใหม่
- `Retry-After` บอกผู้เรียกว่าให้รออีกกี่วินาที เป็นมารยาทที่ควรมีเสมอ เพราะช่วยให้ผู้เรียกที่สุจริตเขียนโปรแกรมรอได้ถูก

<aside>
⚠️

Fixed Window มีจุดอ่อนที่รอยต่อของช่วงเวลา

หากจำกัดไว้ 5 ครั้งต่อนาที ผู้โจมตียิง 5 ครั้งตอนวินาทีที่ 59 แล้วยิงอีก 5 ครั้งตอนวินาทีที่ 61 จะได้ 10 ครั้งภายในสองวินาที

วิธีที่แม่นกว่าคือ Sliding Window หรือ Token Bucket ซึ่งซับซ้อนขึ้นและอยู่นอกขอบเขตของวันนี้

</aside>

**ไฟล์ที่ 2 — `delivery/http/middleware_headers.go` สร้างใหม่**

```go
package http

import "github.com/gofiber/fiber/v2"

// SecurityHeaders ใส่ Header ป้องกันพื้นฐานให้ทุก Response
func SecurityHeaders(c *fiber.Ctx) error {
	// ห้ามเบราว์เซอร์เดาชนิดไฟล์เอง ป้องกันการหลอกให้รันไฟล์ผิดชนิด
	c.Set("X-Content-Type-Options", "nosniff")
	// ห้ามนำหน้าเว็บไปฝังใน iframe ของเว็บอื่น ป้องกัน Clickjacking
	c.Set("X-Frame-Options", "DENY")
	// ไม่ส่ง URL เดิมติดไปกับ Request ที่ออกไปยังเว็บอื่น
	c.Set("Referrer-Policy", "no-referrer")
	return c.Next()
}
```

| Header | ป้องกันอะไร |
| --- | --- |
| `X-Content-Type-Options: nosniff` | ห้ามเบราว์เซอร์เดาชนิดไฟล์เอง ป้องกันการหลอกให้รันไฟล์ผิดชนิด |
| `X-Frame-Options: DENY` | ห้ามนำหน้าเว็บไปฝังใน `iframe` ของเว็บอื่น ป้องกัน Clickjacking |
| `Referrer-Policy: no-referrer` | ไม่ส่ง URL เดิมติดไปกับ Request ที่ออกไปยังเว็บอื่น ป้องกันข้อมูลใน URL รั่ว |

## 3.9 CORS และการติดตั้ง Middleware

### 3.9.1 CORS คืออะไร และไม่ใช่อะไร

CORS เป็นหัวข้อที่สร้างความสับสนมากที่สุดหัวข้อหนึ่ง เพราะคนส่วนใหญ่เข้าใจผิดตั้งแต่ประโยคแรก จึงต้องเริ่มจากการบอกว่ามันไม่ใช่อะไร

<aside>
🚨

**CORS ไม่ใช่ระบบป้องกัน API ของเรา**

CORS เป็นกฎของ **เบราว์เซอร์** ไม่ใช่กฎของเซิร์ฟเวอร์ ผู้ที่บังคับใช้กฎนี้คือเบราว์เซอร์ ไม่ใช่โค้ด Go ของเรา

`curl`, Postman, สคริปต์ Python และโปรแกรมทุกตัวที่ไม่ใช่เบราว์เซอร์ **ไม่สนใจ CORS เลย** ยิงเข้ามาได้ตามปกติ

สังเกตว่าทั้งวันนี้เรายิง `curl` เข้า API โดยไม่เคยติด CORS แม้แต่ครั้งเดียว ทั้งที่ตั้งค่า CORS ไว้แล้ว

**สิ่งที่ปกป้อง API คือ Token และการตรวจสิทธิ์ที่ทำมาทั้งวัน ไม่ใช่ CORS**

</aside>

**แล้ว CORS มีไว้ทำอะไร**

เบราว์เซอร์มีกฎพื้นฐานข้อหนึ่งชื่อ Same-Origin Policy ใจความคือ *โค้ด JavaScript ของเว็บหนึ่ง ห้ามอ่านผลลัพธ์จากอีกเว็บหนึ่ง* กฎนี้มีไว้กันไม่ให้เว็บอันตรายที่เราเผลอเปิด ดึงข้อมูลจากเว็บธนาคารที่เราล็อกอินค้างไว้ในอีกแท็บหนึ่ง

CORS คือช่องทางที่เซิร์ฟเวอร์ใช้บอกเบราว์เซอร์ว่า *เว็บใดบ้างที่ฉันอนุญาตให้อ่านผลลัพธ์ของฉันได้*

**Origin คืออะไร** ประกอบด้วยสามส่วน ต่างกันแม้ส่วนเดียวก็ถือว่าคนละ Origin

| เทียบกับ `http://localhost:5173` | ถือว่า | เพราะ |
| --- | --- | --- |
| `http://localhost:5173/orders` | Origin เดียวกัน | เส้นทางไม่นับ นับแค่ Scheme, Host, Port |
| `http://localhost:3000` | **คนละ** Origin | Port ต่างกัน |
| `https://localhost:5173` | **คนละ** Origin | Scheme ต่างกัน |
| `http://127.0.0.1:5173` | **คนละ** Origin | Host ต่างกัน แม้จะเป็นเครื่องเดียวกันก็ตาม |

<aside>
💡

แถวสุดท้ายคือสาเหตุของอาการ *เมื่อวานยังใช้ได้ วันนี้ใช้ไม่ได้* ที่พบบ่อยที่สุด

`localhost` กับ `127.0.0.1` เป็นคนละ Origin ในสายตาเบราว์เซอร์ ทั้งที่ชี้ไปเครื่องเดียวกัน

</aside>

**สิ่งที่เบราว์เซอร์บล็อกจริง ๆ**

จุดนี้สำคัญและมักเข้าใจผิด สำหรับ Request ธรรมดา เบราว์เซอร์ **ส่งไปถึงเซิร์ฟเวอร์แล้ว** และเซิร์ฟเวอร์ **ทำงานไปแล้วเรียบร้อย** สิ่งที่เบราว์เซอร์บล็อกคือการที่โค้ด JavaScript จะอ่านผลลัพธ์กลับมา

แปลว่าหากเส้นทางนั้นลบข้อมูลทิ้ง ข้อมูลก็ถูกลบไปแล้ว เพียงแต่โค้ดหน้าเว็บอ่านคำตอบไม่ได้เท่านั้น **นี่คือเหตุผลอีกข้อที่ CORS ไม่นับเป็นระบบป้องกัน**

**Preflight คำถามนำก่อนส่งของจริง**

Request ที่ไม่ธรรมดา เช่น ใช้ `PATCH` หรือ `DELETE` หรือแนบ Header ชื่อ `Authorization` เบราว์เซอร์จะไม่ส่งของจริงทันที แต่ยิง `OPTIONS` ไปถามก่อนหนึ่งครั้ง

```mermaid
sequenceDiagram
    participant JS as โค้ดหน้าเว็บ
    participant B as เบราว์เซอร์
    participant API as Fiber API

    JS->>B: fetch POST /api/orders พร้อม Authorization
    Note over B: เป็น Request ที่ไม่ธรรมดา ต้องถามก่อน
    B->>API: OPTIONS /api/orders<br/>Origin: http://localhost:5173<br/>Access-Control-Request-Headers: Authorization
    API-->>B: 204<br/>Access-Control-Allow-Origin: http://localhost:5173<br/>Access-Control-Allow-Headers: Content-Type,Authorization
    Note over B: ต้นทางและ Header อยู่ในรายการอนุญาต
    B->>API: POST /api/orders พร้อม Authorization
    API-->>B: 201 Created
    B-->>JS: อ่านผลลัพธ์ได้
```

เมื่อเข้าใจภาพนี้แล้ว ค่าที่ตั้งไว้ในขั้นที่ 12 จะอ่านออกทันทีว่าแต่ละบรรทัดตอบคำถามข้อใดของเบราว์เซอร์

```go
app.Use(cors.New(cors.Config{
	AllowOrigins: envOr("CORS_ORIGINS", "http://localhost:5173"),
	AllowMethods: "GET,POST,PATCH,DELETE",
	AllowHeaders: "Content-Type,Authorization",
}))
```

| บรรทัด | ตอบคำถามว่า |
| --- | --- |
| `AllowOrigins` | เว็บใดบ้างที่อ่านผลลัพธ์ของเราได้ ระบุเป็นรายชื่อ ไม่ใช้ `*` |
| `AllowMethods` | ใช้ Method ใดได้บ้าง ต้องมี `PATCH` และ `DELETE` เพราะสองตัวนี้ทำให้เกิด Preflight |
| `AllowHeaders` | แนบ Header ใดได้บ้าง ต้องมี `Authorization` เพราะเราส่ง Token ไปกับทุก Request |

<aside>
⚠️

**ห้ามใช้ `AllowOrigins: "*"` คู่กับการส่ง Credential**

เบราว์เซอร์ปฏิเสธการจับคู่นี้ตรง ๆ ตามมาตรฐาน เมื่อใดที่เปิด `AllowCredentials` แล้วยังใส่ `*` ไว้ เบราว์เซอร์จะบล็อกทุก Request

และแม้ในกรณีที่ใช้ได้ `*` ก็แปลว่าเว็บใดบนอินเทอร์เน็ตก็อ่านผลลัพธ์ของเราได้ ซึ่งแทบไม่เคยเป็นสิ่งที่ต้องการ

ระบุรายชื่อต้นทางให้ชัดเจนเสมอ และเก็บไว้ใน `.env` เพราะต้นทางของเครื่องพัฒนากับของจริงไม่เหมือนกัน

</aside>

**ตารางวินิจฉัยเวลาเจอข้อความ CORS error**

| ข้อความที่เห็นใน Console | ความหมายและวิธีแก้ |
| --- | --- |
| `No 'Access-Control-Allow-Origin' header` | ต้นทางของหน้าเว็บไม่อยู่ใน `CORS_ORIGINS` ให้เติมเข้าไปแล้วรีสตาร์ทเซิร์ฟเวอร์ |
| `Method PATCH is not allowed` | ยังไม่ได้ใส่ `PATCH` ใน `AllowMethods` |
| `Request header field authorization is not allowed` | ยังไม่ได้ใส่ `Authorization` ใน `AllowHeaders` |
| ยิงด้วย `curl` ผ่าน แต่หน้าเว็บไม่ผ่าน | ปกติ เพราะ `curl` ไม่บังคับใช้ CORS อาการนี้ยืนยันว่าเป็นเรื่อง CORS จริง ไม่ใช่ API พัง |
| ตอบ 401 แต่ Console ฟ้อง CORS | ดู Status Code จริงในแท็บ Network ก่อนเสมอ ปัญหาที่แท้จริงมักเป็นเรื่อง Token ไม่ใช่ CORS |

#### ขั้นที่ 12 — ติดตั้ง Middleware ทั้งหมด พร้อม CORS, Timeout และขนาด Body

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go`

**2. แก้ไขตามรายการด้านล่าง** หรือจะลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**3. รัน** `go run .`

ขั้นนี้ติดตั้ง Middleware ทั้งสองตัวของขั้นที่แล้วพร้อมกัน แล้วทดสอบทั้งคู่ในรอบเดียว

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไข**

- สร้างตัวนับสองชุด ชุดกว้างสำหรับทุกเส้นทาง และชุดเข้มเฉพาะหน้าเข้าสู่ระบบ

```go
// ตัวนับสองชุด ชุดกว้างสำหรับทุกเส้นทาง และชุดเข้มเฉพาะหน้าเข้าสู่ระบบ
general := http.NewRateLimiter(100, time.Minute)
loginer := http.NewRateLimiter(5, time.Minute)
app.Use(general.Limit)
```

แล้วผูกตัวนับชุดเข้มไว้เฉพาะเส้นทางเข้าสู่ระบบ

```go
api.Post("/auth/login", loginer.Limit, h.Login)
```

<aside>
🚨

ตัวนับนี้เก็บอยู่ใน `map` ในหน่วยความจำของ Process จึงเป็น **Stateful**

เมื่อรันสองสำเนาพร้อมกันหลัง Load Balancer แต่ละสำเนาจะนับแยกกัน ผู้โจมตีจึงยิงได้ 5 ครั้งต่อสำเนา รวมเป็น 10 ครั้ง

นี่คือหัวข้อ 2.5.1 ของบทที่แล้วกลับมาอีกครั้ง คราวนี้พร้อมผลเสียที่จับต้องได้

งานจริงจึงต้องย้ายตัวนับออกไปไว้ที่ Redis ซึ่งเป็นที่เก็บส่วนกลางที่ทุกสำเนาเห็นตรงกัน

</aside>

**สิ่งที่แก้ไข**

- ตั้งค่า `fiber.Config` จำกัดขนาด Body และตัดการเชื่อมต่อที่ค้างนาน
- ใส่ `SecurityHeaders` และ CORS ให้ครอบทุก Request

```go
app := fiber.New(fiber.Config{
	// จำกัดขนาด Body ป้องกันการส่งไฟล์ใหญ่มากเข้ามาจนหน่วยความจำหมด
	BodyLimit: 1 * 1024 * 1024,
	// ตัดการเชื่อมต่อที่ค้างนานเกินไป
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
})

app.Use(http.SecurityHeaders)

// อนุญาตเฉพาะต้นทางที่ระบุไว้ ห้ามใช้ * คู่กับการส่ง Credential
app.Use(cors.New(cors.Config{
	AllowOrigins: envOr("CORS_ORIGINS", "http://localhost:5173"),
	AllowMethods: "GET,POST,PATCH,DELETE",
	AllowHeaders: "Content-Type,Authorization",
}))
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"context"
    	"log"
    	"os"
    	"strings"
    	"time"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/repository/postgres"
    	"food-api/security"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    	"github.com/gofiber/fiber/v2/middleware/cors"
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
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึงทุกชั้นพร้อมกัน และเป็นที่เดียวที่เลือกว่า
    // จะใช้ Adapter ตัวใด ทั้งที่เก็บข้อมูล วิธี Hash รหัสผ่าน และรูปแบบ Token
    func main() {
    	_ = godotenv.Load()
    
    	// ค่าลับต้องมาจาก Environment เสมอ และต้องหยุดทำงานทันทีเมื่อไม่มี
    	// ห้ามใส่ค่าสำรองไว้ในโค้ด เพราะจะกลายเป็นรหัสลับที่ทุกคนบนโลกรู้
    	secret := os.Getenv("JWT_SECRET")
    	if len(strings.TrimSpace(secret)) < 32 {
    		log.Fatal("ต้องกำหนด JWT_SECRET ใน .env และต้องยาวอย่างน้อย 32 ตัวอักษร")
    	}
    
    	ttl, err := time.ParseDuration(envOr("JWT_TTL", "1h"))
    	if err != nil {
    		log.Fatal("JWT_TTL ไม่ถูกต้อง: ", err)
    	}
    
    	// 1. ที่เก็บข้อมูล เลือกจาก Environment เช่นเดิม
    	var orders usecase.OrderRepository
    	var users usecase.UserRepository
    
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
    		orders = postgres.NewOrderRepo(pool)
    		users = postgres.NewUserRepo(pool)
    	} else {
    		log.Println("ใช้ Memory")
    		orders = memory.NewOrderRepo(seedMenu())
    		users = memory.NewUserRepo()
    	}
    
    	// 2. Adapter ด้านความปลอดภัย
    	//    ชั้น Use Case รู้จักเพียง Interface ไม่รู้จักชื่อ bcrypt และ JWT
    	hasher := security.NewBcryptHasher(0)
    
    	tokens := security.NewJWTIssuer(secret, ttl)
    
    	// 3. กฎ
    	uc := usecase.NewOrderUseCase(orders)
    	auth := usecase.NewAuthUseCase(users, hasher, tokens)
    
    	// บัญชีพนักงานคนแรกมาจากผู้ดูแลระบบ ไม่ใช่จากการสมัครเอง
    	// เส้นทาง /auth/register สาธารณะสร้างได้เฉพาะลูกค้าเท่านั้น
    	if email, pw := os.Getenv("STAFF_EMAIL"), os.Getenv("STAFF_PASSWORD"); email != "" && pw != "" {
    		if _, err := auth.CreateStaff(email, pw); err != nil {
    			log.Println("ข้ามการสร้างบัญชีพนักงานเริ่มต้น:", err)
    		} else {
    			log.Println("สร้างบัญชีพนักงานเริ่มต้นแล้ว:", email)
    		}
    	}
    
    	// 4. ชั้นขนส่ง
    	h := http.NewHandler(uc, auth)
    
    	app := fiber.New(fiber.Config{
    		// จำกัดขนาด Body ป้องกันการส่งไฟล์ใหญ่มากเข้ามาจนหน่วยความจำหมด
    		BodyLimit: 1 * 1024 * 1024,
    		// ตัดการเชื่อมต่อที่ค้างนานเกินไป
    		ReadTimeout:  10 * time.Second,
    		WriteTimeout: 10 * time.Second,
    	})
    
    	app.Use(http.SecurityHeaders)
    
    	// อนุญาตเฉพาะต้นทางที่ระบุไว้ ห้ามใช้ * คู่กับการส่ง Credential
    	app.Use(cors.New(cors.Config{
    		AllowOrigins: envOr("CORS_ORIGINS", "http://localhost:5173"),
    		AllowMethods: "GET,POST,PATCH,DELETE",
    		AllowHeaders: "Content-Type,Authorization",
    	}))
    
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	// ตัวนับสองชุด ชุดกว้างสำหรับทุกเส้นทาง และชุดเข้มเฉพาะหน้าเข้าสู่ระบบ
    	general := http.NewRateLimiter(100, time.Minute)
    	loginer := http.NewRateLimiter(5, time.Minute)
    	app.Use(general.Limit)
    
    	api := app.Group("/api")
    
    	// ── เปิดให้ทุกคน ──
    	api.Post("/auth/register", h.Register)
    	api.Post("/auth/login", loginer.Limit, h.Login)
    	api.Get("/menu", h.ListMenu)
    
    	// ── ต้องเข้าสู่ระบบก่อน ──
    	me := api.Group("", h.RequireAuth)
    	me.Get("/auth/me", h.Me)
    	me.Post("/orders", h.PlaceOrder)
    	me.Get("/orders/:id", h.GetOrder)
    	me.Get("/orders/:id/total", h.OrderTotal)
    	me.Post("/orders/:id/cancel", h.CancelOrder)
    
    	// ── เฉพาะพนักงานร้าน ──
    	staff := api.Group("", h.RequireAuth, h.RequireStaff)
    	staff.Patch("/orders/:id/status", h.SetStatus)
    
    	log.Fatal(app.Listen(":3000"))
    }
    
    func envOr(key, fallback string) string {
    	if v := os.Getenv(key); v != "" {
    		return v
    	}
    	return fallback
    }
    ```
    

<aside>
🚨

**ห้ามตั้ง `AllowOrigins` เป็น `*` เมื่อ API ต้องรับ Token**

การเปิดกว้างแบบนั้น เท่ากับอนุญาตให้เว็บใดก็ได้บนอินเทอร์เน็ตเรียก API ของเราจากเบราว์เซอร์ของผู้ใช้ที่ล็อกอินอยู่

ให้ระบุรายชื่อต้นทางที่อนุญาตให้ชัดเจนเสมอ และอ่านค่ามาจาก `.env` เพื่อให้แต่ละที่ติดตั้งกำหนดเองได้

</aside>

- `BodyLimit` และ Timeout เป็นการป้องกันขั้นพื้นฐานที่มักถูกลืม เพราะไม่มีใครสังเกตจนกว่าจะโดนโจมตีจริง

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
for i in 1 2 3 4 5 6 7; do
  curl -s -o /dev/null -w "ครั้งที่ $i -> %{http_code}\n" \
    -X POST localhost:3000/api/auth/login -H 'Content-Type: application/json' \
    -d '{"email":"alice@example.com","password":"password123"}'
done
```

```bash
ครั้งที่ 1 -> 200
ครั้งที่ 2 -> 200
ครั้งที่ 3 -> 200
ครั้งที่ 4 -> 200
ครั้งที่ 5 -> 200
ครั้งที่ 6 -> 429
ครั้งที่ 7 -> 429
```

คำสั่งที่ 2

```bash
curl -i -X POST localhost:3000/api/auth/login -H 'Content-Type: application/json' -d '{"email":"alice@example.com","password":"password123"}'
```

```bash
HTTP/1.1 429 Too Many Requests
X-Ratelimit-Limit: 5
X-Ratelimit-Remaining: 0
Retry-After: 60

{"error":"too many requests"}
```

คำสั่งที่ 3

```bash
curl -I localhost:3000/api/menu
```

```bash
HTTP/1.1 200 OK
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: no-referrer
X-Ratelimit-Limit: 100
X-Ratelimit-Remaining: 99
```

## 3.10 พิสูจน์ว่ากฎทำงานจริง

ขั้นนี้ไม่ต้องเขียนโค้ดเพิ่ม เป็นการทดลองยิงจริงเพื่อดูว่ากฎที่เขียนไว้ทำงานตามที่ตั้งใจหรือไม่

#### ขั้นที่ 13 — พิสูจน์ว่ากฎทำงานจริง และ Token ปลอมใช้ไม่ได้

<aside>
📄

**ขั้นนี้ไม่ต้องเขียนโค้ดเพิ่ม เป็นการยิงจริงเพื่อพิสูจน์กฎที่เขียนไว้ทั้งหมด**

**1. เตรียมข้อมูล** สมัครลูกค้าสองคน ตามคำสั่งด้านล่าง ส่วนบัญชีพนักงานถูกสร้างให้แล้วตอนเปิดเซิร์ฟเวอร์

**2. ชุดที่หนึ่ง** พิสูจน์ว่าลูกค้าคนหนึ่ง แตะออเดอร์ของลูกค้าอีกคนไม่ได้

**3. ชุดที่สอง** พิสูจน์ว่า Token ปลอมและ Token หมดอายุ ใช้ไม่ได้

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ชุดที่หนึ่ง — ลูกค้าคนหนึ่ง แตะออเดอร์ของลูกค้าอีกคนไม่ได้**

**เตรียมข้อมูล** สมัครลูกค้าสองคน แล้วเก็บ Token ของทั้งสามบัญชีไว้

```bash
B=localhost:3000/api

# สมัครลูกค้าสองคน เส้นทางสาธารณะสร้างได้เฉพาะลูกค้าเท่านั้น
curl -s -X POST $B/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123"}'
curl -s -X POST $B/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"bob@example.com","password":"password123"}'

# บัญชีพนักงานไม่ได้มาจากการสมัคร แต่ถูกสร้างตอนเปิดเซิร์ฟเวอร์
# จากค่า STAFF_EMAIL และ STAFF_PASSWORD ใน .env ของขั้นที่ 6

login() { curl -s -X POST $B/auth/login -H 'Content-Type: application/json' \
  -d "{\"email\":\"$1\",\"password\":\"password123\"}" | jq -r .token; }

ALICE=$(login alice@example.com)
BOB=$(login bob@example.com)
STAFF=$(login staff@example.com)
```

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
curl -s -X POST localhost:3000/api/orders -H "Authorization: Bearer $ALICE" -H 'Content-Type: application/json' -d '{"menu_id":"m1","qty":2}'
```

```bash
{"id":"o1","status":"pending","total":160,"lines":[{"menu_id":"m1","name":"Pad Thai","price":80,"qty":2}]}
```

คำสั่งที่ 2

```bash
curl -s -o /dev/null -w "%{http_code}\n" localhost:3000/api/orders/o1 -H "Authorization: Bearer $ALICE" 
```

```bash
200
```

คำสั่งที่ 3

```bash
curl -s -i localhost:3000/api/orders/o1 -H "Authorization: Bearer $BOB" 
```

```bash
HTTP/1.1 404 Not Found

{"error":"not found"}
```

คำสั่งที่ 4

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X POST localhost:3000/api/orders/o1/cancel -H "Authorization: Bearer $BOB" 
```

```bash
404
```

คำสั่งที่ 5

```bash
curl -s -i -X PATCH localhost:3000/api/orders/o1/status -H "Authorization: Bearer $ALICE" -H 'Content-Type: application/json' -d '{"status":"confirmed"}'
```

```bash
HTTP/1.1 403 Forbidden

{"error":"forbidden"}
```

คำสั่งที่ 6

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X PATCH localhost:3000/api/orders/o1/status -H "Authorization: Bearer $STAFF" -H 'Content-Type: application/json' -d '{"status":"confirmed"}'
```

```bash
204
```

คำสั่งที่ 7

```bash
curl -s -o /dev/null -w "%{http_code}\n" -X POST localhost:3000/api/orders/o1/cancel -H "Authorization: Bearer $ALICE" 
```

```bash
204
```

<aside>
💡

สังเกตความต่างของคำสั่งที่ 3 กับคำสั่งที่ 5 ให้ดี

**คำสั่งที่ 3** bob ขอดูออเดอร์ของ alice ได้ `404` เพราะระบบไม่ยอมยืนยันแม้แต่ว่าออเดอร์ใบนี้มีอยู่จริง

**คำสั่งที่ 5** alice ขอเปลี่ยนสถานะออเดอร์ *ของตนเอง* ได้ `403` เพราะระบบยอมรับว่าเห็นออเดอร์ใบนี้ได้ เพียงแต่ทำสิ่งนี้ไม่ได้

การเลือกใช้ `404` หรือ `403` จึงขึ้นกับว่า **ผู้เรียกควรรู้หรือไม่ว่าสิ่งนั้นมีอยู่**

</aside>

- ทั้งเจ็ดคำสั่งไม่มีข้อใดที่ตัดสินใจอยู่ในชั้น `delivery` เลย ทุกคำตอบมาจากกฎในชั้น `domain` ที่เขียนไว้ตั้งแต่ขั้นที่ 2

**ชุดที่สอง — Token ปลอมและ Token หมดอายุ ใช้ไม่ได้**

**ผลลัพธ์** — เปิด Terminal อีกหน้าต่างหนึ่ง โดยปล่อยให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกใช้คำสั่งด้านล่างทีละคำสั่ง

คำสั่งที่ 1

```bash
echo $ALICE | cut -d. -f2 | base64 -d
```

```bash
{"email":"alice@example.com","exp":1789014380,"iat":1789010780,"role":"customer","sub":"u1
```

คำสั่งที่ 2

```bash
BAD="${ALICE%?}X"
curl -s -i localhost:3000/api/auth/me -H "Authorization: Bearer $BAD" 
```

```bash
HTTP/1.1 401 Unauthorized

{"error":"unauthorized"}
```

คำสั่งที่ 3

```bash
# หลังแก้ JWT_TTL เป็น 10s เปิดเซิร์ฟเวอร์ใหม่ ขอ Token ใหม่ แล้วรอให้เกินสิบวินาที
sleep 15
curl -s -o /dev/null -w "%{http_code}\n" localhost:3000/api/auth/me -H "Authorization: Bearer $SHORT" 
```

```bash
401
```

<aside>
🚨

คำสั่งที่ 1 พิสูจน์ว่า **Payload อ่านได้ด้วยตาเปล่า** โดยไม่ต้องมีกุญแจใด ๆ

หากเผลอใส่รหัสผ่านหรือเลขบัตรประชาชนลงไป ผู้ที่ได้ Token ไปจะอ่านได้ทันที

สิ่งที่ลายเซ็นรับประกันคือ *ไม่ถูกแก้* ไม่ใช่ *เป็นความลับ*

</aside>

- ข้อความที่ได้จากคำสั่งที่ 1 อาจขาดตัวอักษรท้ายไปหนึ่งถึงสองตัว เพราะความยาวของ Base64 ไม่ลงตัวพอดีสี่ตัวอักษร ไม่ใช่ข้อผิดพลาด
- คำสั่งที่ 2 แก้ลายเซ็นเพียงตัวอักษรเดียว ระบบก็ปฏิเสธทันที เพราะลายเซ็นคำนวณจากทั้ง Header และ Payload รวมกับกุญแจลับ

<aside>
💡

ลองแก้ `role` ใน Payload เป็น `staff` แล้วประกอบ Token ขึ้นใหม่ดู จะพบว่าใช้ไม่ได้เช่นกัน

เพราะเมื่อ Payload เปลี่ยน ลายเซ็นเดิมจะไม่ตรงอีกต่อไป และผู้ที่ไม่มี `JWT_SECRET` สร้างลายเซ็นใหม่ให้ตรงไม่ได้

นี่คือเหตุผลทั้งหมดที่ `JWT_SECRET` ต้องเป็นความลับ และต้องยาวพอที่จะเดาไม่ได้

</aside>

## 3.11 ตัวอย่างที่สอง — Session แบบ Stateful

ถึงตรงนี้ระบบยืนยันตัวตนด้วย JWT ทำงานครบแล้ว หัวข้อนี้จะสร้างวิธีที่สองขึ้นมาคู่กัน แล้วสลับไปมาได้ด้วยการแก้ `.env` เพียงบรรทัดเดียว

เหตุผลที่ทำได้ ย้อนกลับไปที่ขั้นที่ 3 ซึ่งชั้น Use Case ประกาศไว้ว่าตนต้องการอะไร โดยไม่ระบุว่าใครจะเป็นผู้ทำให้

```go
type TokenIssuer interface {
	Issue(u domain.User) (string, error)
	Parse(token string) (domain.User, error)
}
```

<aside>
💡

**Interface ตัวนี้ไม่มีคำว่า JWT อยู่เลย**

จึงไม่ได้ผูกกับ JWT มาตั้งแต่ต้น สิ่งที่เขียนไว้คือ *แปลงผู้ใช้เป็นข้อความ และแปลงข้อความกลับเป็นผู้ใช้* เท่านั้น

Session ก็ทำสองอย่างนี้ได้เหมือนกัน เพียงแต่ทำด้วยวิธีที่ต่างออกไปโดยสิ้นเชิง

นี่คือผลตอบแทนของ Dependency Inversion ที่ลงทุนไว้ตั้งแต่ Day 2 หัวข้อ 2.5.6 ตอนสลับ Memory เป็น Postgres

</aside>

**ความต่างอยู่ที่ว่า ใครเป็นคนจำ**

|  | JWT | Session |
| --- | --- | --- |
| ข้อมูลตัวตนอยู่ที่ไหน | อยู่ในตัว Token ที่ผู้ใช้ถือไป | อยู่ในฐานข้อมูลฝั่งเซิร์ฟเวอร์ |
| ตัวข้อความที่ส่งไปมา | ยาว อ่านออกด้วยตาเปล่า มีข้อมูลจริงอยู่ข้างใน | สั้น สุ่มล้วน ไม่มีความหมายในตัวเอง |
| เซิร์ฟเวอร์ตรวจอย่างไร | ตรวจลายเซ็นด้วยกุญแจลับ ไม่ต้องถามใคร | เอาไปค้นในตาราง `sessions` |
| ถ้าฐานข้อมูลล่ม | ยังตรวจ Token ได้ | ตรวจไม่ได้เลย |

สังเกตแถวสุดท้าย นั่นคือราคาที่ต้องจ่าย และแถวถัดไปคือสิ่งที่ได้กลับมา

#### ขั้นที่ 14 — Adapter ตัวที่ 3 — ออกและตรวจ Session

<aside>
📄

**1. สร้างไฟล์** `security/session.go` แล้ววางโค้ดด้านล่างลงไป

**2. รัน** `go build ./...`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

```go
package security

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"food-api/domain"
)

// Session คือสิ่งที่เซิร์ฟเวอร์จำไว้ฝั่งตนเอง
// ตรงข้ามกับ JWT ที่ไม่จำอะไรเลย แล้วฝากข้อมูลทั้งหมดไว้กับผู้ใช้
type Session struct {
	UserID    string
	Email     string
	Role      string
	ExpiresAt time.Time
}

// SessionStore ประกาศไว้ที่นี่ เพราะที่นี่คือผู้ใช้งาน
// หลักการเดียวกับ usecase.UserRepository ที่ประกาศไว้ในชั้นที่เรียกใช้ ไม่ใช่ชั้นที่ให้บริการ
type SessionStore interface {
	Save(id string, s Session) error
	Find(id string) (Session, error)
	Delete(id string) error
}

// SessionIssuer ทำตาม Interface usecase.TokenIssuer ตัวเดียวกันกับที่ JWTIssuer ทำ
//
// เพราะทำตาม Interface เดียวกัน ชั้น Use Case จึงสลับมาใช้ตัวนี้ได้
// โดยไม่ต้องแก้โค้ดในชั้น Domain ชั้น Use Case และ Handler แม้แต่บรรทัดเดียว
type SessionIssuer struct {
	store SessionStore
	ttl   time.Duration
}

func NewSessionIssuer(s SessionStore, ttl time.Duration) SessionIssuer {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	return SessionIssuer{store: s, ttl: ttl}
}

func (s SessionIssuer) Issue(u domain.User) (string, error) {
	// Session ID ต้องเดาไม่ได้ จึงต้องสุ่มด้วย crypto/rand เท่านั้น
	// ห้ามใช้ math/rand เด็ดขาด เพราะผู้โจมตีคาดเดาค่าถัดไปได้
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	id := base64.RawURLEncoding.EncodeToString(b)

	// ต่างจาก JWT ตรงที่ตัว ID ไม่ได้บอกอะไรเลย เป็นเพียงกุญแจไว้เปิดหาข้อมูลจริง
	err := s.store.Save(id, Session{
		UserID:    u.ID,
		Email:     u.Email,
		Role:      string(u.Role),
		ExpiresAt: time.Now().Add(s.ttl),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s SessionIssuer) Parse(token string) (domain.User, error) {
	// ทุก Request ต้องถามที่เก็บข้อมูลหนึ่งครั้ง
	// นี่คือราคาที่จ่ายไป เพื่อแลกกับความสามารถในการยกเลิกได้ทันที
	sess, err := s.store.Find(token)
	if err != nil {
		return domain.User{}, domain.ErrUnauthorized
	}

	if time.Now().After(sess.ExpiresAt) {
		_ = s.store.Delete(token)
		return domain.User{}, domain.ErrUnauthorized
	}

	return domain.User{
		ID:    sess.UserID,
		Email: sess.Email,
		Role:  domain.Role(sess.Role),
	}, nil
}

// Revoke ทำให้ Token ตายทันที เป็นสิ่งที่ JWTIssuer ทำไม่ได้
// การมี Method นี้ทำให้ SessionIssuer ทำตาม Interface usecase.TokenRevoker ไปด้วย
func (s SessionIssuer) Revoke(token string) error {
	return s.store.Delete(token)
}
```

**ผลลัพธ์**

```bash
go build ./...
```

- ไฟล์นี้ทำตาม Interface `usecase.TokenIssuer` **ตัวเดียวกัน** กับที่ `security/jwt.go` ทำ ไม่มีการแก้ Interface แม้แต่ตัวอักษรเดียว
- `SessionStore` ถูกประกาศไว้ในไฟล์นี้ เพราะไฟล์นี้คือ *ผู้ใช้งาน* หลักการเดียวกับที่ `usecase` ประกาศ `UserRepository` ไว้ในขั้นที่ 3

<aside>
🚨

**Session ID ต้องมาจาก `crypto/rand` เท่านั้น**

`math/rand` สร้างเลขสุ่มที่คาดเดาได้ หากผู้โจมตีรู้ค่าที่ออกมาไม่กี่ค่า จะคำนวณค่าถัดไปได้ทั้งหมด แล้วสวมรอยเป็นผู้ใช้คนอื่นได้ทันที

ความแตกต่างนี้ไม่มีทางเห็นจากการทดสอบ เพราะทั้งสองแบบ *ดูสุ่มเหมือนกัน* ในสายตาคน

กฎที่จำง่ายคือ ค่าใดที่ใช้แทนตัวตนหรือใช้เป็นความลับ ให้ใช้ `crypto/rand` เสมอ

</aside>

- `Parse` ต้องอ่านที่เก็บข้อมูลหนึ่งครั้งทุก Request ซึ่ง `JWTIssuer` ไม่ต้องทำ นี่คือต้นทุนที่จ่ายเพื่อแลกกับความสามารถในขั้นที่ 16
- `Revoke` เป็น Method ที่ `JWTIssuer` ไม่มี ความต่างข้อนี้จะถูกใช้ในขั้นที่ 16

#### ขั้นที่ 15 — ที่เก็บ Session ทั้งแบบ Memory และ Postgres พร้อม Migration ที่ 6

<aside>
📄

**1. สร้างไฟล์ที่ 1** `repository/memory/session.go` แล้ววางโค้ดด้านล่างลงไป

**2. สร้างไฟล์ที่ 2** `repository/postgres/session.go` แล้ววางโค้ดด้านล่างลงไป

**3. สร้างไฟล์ที่ 3** `migrations/006_sessions.sql` แล้ววางโค้ดด้านล่างลงไป

**4. รัน** `go run ./cmd/migrate`

ขั้นนี้มี **3 ไฟล์** สร้างให้ครบทุกไฟล์ก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

ไฟล์ `repository/memory/session.go`

```go
package memory

import (
	"sync"

	"food-api/domain"
	"food-api/security"
)

// SessionRepo ทำตาม Interface security.SessionStore
//
// เก็บไว้ในหน่วยความจำของโปรเซส จึงหายทั้งหมดเมื่อรีสตาร์ท
// และถ้าเพิ่มสำเนาเซิร์ฟเวอร์เป็นสองตัว แต่ละตัวจะจำคนละชุด นี่คือข้อเสียของ Stateful
type SessionRepo struct {
	mu       *sync.RWMutex
	sessions map[string]security.Session
}

func NewSessionRepo() SessionRepo {
	return SessionRepo{mu: &sync.RWMutex{}, sessions: map[string]security.Session{}}
}

func (r SessionRepo) Save(id string, s security.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[id] = s
	return nil
}

func (r SessionRepo) Find(id string) (security.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sessions[id]
	if !ok {
		return security.Session{}, domain.ErrNotFound
	}
	return s, nil
}

func (r SessionRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, id)
	return nil
}
```

ไฟล์ `repository/postgres/session.go`

```go
package postgres

import (
	"context"

	"food-api/domain"
	"food-api/security"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SessionRepo ทำตาม Interface security.SessionStore
//
// เก็บลงฐานข้อมูล จึงไม่หายเมื่อรีสตาร์ท และสำเนาเซิร์ฟเวอร์ทุกตัวเห็นชุดเดียวกัน
// ระบบจริงมักใช้ Redis แทน เพราะอ่านเร็วกว่าและตั้งเวลาให้หมดอายุเองได้
type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) SessionRepo {
	return SessionRepo{pool: pool}
}

func (r SessionRepo) Save(id string, s security.Session) error {
	// ใช้ $1 $2 เช่นเดิม ค่าที่ส่งเข้ามาไม่มีทางกลายเป็นคำสั่ง SQL
	_, err := r.pool.Exec(context.Background(),
		`INSERT INTO sessions (id, user_id, email, role, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		id, s.UserID, s.Email, s.Role, s.ExpiresAt)
	return err
}

func (r SessionRepo) Find(id string) (security.Session, error) {
	var s security.Session
	err := r.pool.QueryRow(context.Background(),
		`SELECT user_id, email, role, expires_at FROM sessions WHERE id = $1`, id).
		Scan(&s.UserID, &s.Email, &s.Role, &s.ExpiresAt)
	if err != nil {
		return security.Session{}, domain.ErrNotFound
	}
	return s, nil
}

func (r SessionRepo) Delete(id string) error {
	_, err := r.pool.Exec(context.Background(), `DELETE FROM sessions WHERE id = $1`, id)
	return err
}
```

ไฟล์ `migrations/006_sessions.sql`

```sql
-- Migration ที่ 6 ตารางเก็บ Session สำหรับโหมด AUTH_MODE=session
--
-- ตารางนี้คือ "ความจำ" ของเซิร์ฟเวอร์ ซึ่งโหมด JWT ไม่ต้องมี
-- การมีตารางนี้คือราคาที่จ่ายไป เพื่อแลกกับการยกเลิก Token ได้ทันที

CREATE TABLE IF NOT EXISTS sessions (
    id         TEXT PRIMARY KEY,
    user_id    TEXT        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    email      TEXT        NOT NULL,
    role       TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

-- ON DELETE CASCADE ด้านบนทำให้การลบผู้ใช้หนึ่งคน
-- ลบทุก Session ของคนนั้นทิ้งไปด้วย เท่ากับเตะออกจากระบบทุกเครื่องพร้อมกัน

-- ใช้ตอนกวาดล้าง Session ที่หมดอายุทิ้ง
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);
```

**ผลลัพธ์**

```bash
go run ./cmd/migrate
```

```bash
ใช้งาน 006_sessions แล้ว
ใช้งาน migration ใหม่ทั้งหมด 1 ไฟล์
```

- โครงสร้างเหมือน `user.go` ของขั้นที่ 8 ทุกประการ คือมีสองตัวเลือกที่ทำตาม Interface เดียวกัน แล้วให้ `main.go` เป็นผู้เลือก
- `ON DELETE CASCADE` ใน Migration ทำให้การลบผู้ใช้หนึ่งคน ลบทุก Session ของคนนั้นทิ้งไปด้วย เท่ากับเตะออกจากระบบทุกเครื่องพร้อมกัน
- Index บน `expires_at` มีไว้ให้คำสั่งกวาดล้าง Session ที่หมดอายุ ด้วยเหตุผลเดียวกับหัวข้อ 3.7.1

<aside>
💡

ระบบจริงนิยมเก็บ Session ไว้ใน **Redis** มากกว่า Postgres เพราะอ่านเร็วกว่ามาก และตั้งเวลาให้ข้อมูลหายเองได้

แต่รูปแบบของโค้ดเหมือนกันทุกประการ คือเขียน Adapter อีกตัวที่ทำตาม `security.SessionStore` แล้วสลับใน `main.go`

การเปลี่ยนไปใช้ Redis จึงไม่กระทบชั้นอื่นเลย เช่นเดียวกับที่กำลังทำอยู่นี้

</aside>

#### ขั้นที่ 16 — ออกจากระบบ — ความสามารถที่ JWT ทำไม่ได้

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม 4 ไฟล์ ไม่มีไฟล์ใหม่**

**1. เปิดไฟล์** `usecase/auth.go` เพิ่ม Interface และ Method ตามด้านล่าง

**2. เปิดไฟล์** `delivery/http/middleware_auth.go` เก็บ Token ดิบไว้ใน Context

**3. เปิดไฟล์** `delivery/http/auth.go` เพิ่ม Handler

**4. เปิดไฟล์** `delivery/http/error.go` เพิ่มการแปลง Error เป็น 501

**5. รัน** `go build ./...`

ขั้นนี้มี **4 ไฟล์** แต่ละไฟล์แก้ไม่เกินสิบบรรทัด แก้ให้ครบก่อนรันคำสั่ง

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**ไฟล์ที่ 1 — `usecase/auth.go` สิ่งที่แก้ไข**

- ประกาศ Interface เสริมชื่อ `TokenRevoker` ไว้ข้าง ๆ `TokenIssuer`
- เพิ่ม Method `Logout` ซึ่งถาม Adapter ตัวจริงว่ายกเลิกได้หรือไม่

```go
// TokenRevoker เป็นความสามารถเสริม ไม่ใช่ TokenIssuer ทุกตัวที่ทำได้
// จึงแยกออกมาเป็น Interface ต่างหาก แล้วค่อยถามตอนใช้งานว่าตัวจริงทำได้หรือไม่
// รูปแบบนี้พบได้ทั่วไปในไลบรารีมาตรฐานของ Go เช่น io.ReaderFrom
type TokenRevoker interface {
	Revoke(token string) error
}

func (a AuthUseCase) Logout(token string) error {
	r, ok := a.tokens.(TokenRevoker)
	if !ok {
		return domain.ErrRevokeUnsupported
	}
	return r.Revoke(token)
}
```

<aside>
💡

บรรทัด `r, ok := a.tokens.(TokenRevoker)` คือการถามว่า *ของจริงที่เสียบเข้ามา ทำเรื่องนี้ได้ไหม*

`SessionIssuer` มี Method `Revoke` จึงตอบว่าได้  ส่วน `JWTIssuer` ไม่มี จึงตอบว่าไม่ได้

วิธีนี้ทำให้ไม่ต้องยัด `Revoke` เข้าไปใน `TokenIssuer` แล้วบังคับให้ `JWTIssuer` เขียน Method ที่ทำอะไรไม่ได้จริง

หลักการนี้ชื่อ Interface Segregation คือ *อย่าบังคับให้ใครต้องทำตามสัญญาที่ตนไม่ได้ใช้*

</aside>

**ไฟล์ที่ 2 — `delivery/http/middleware_auth.go` สิ่งที่แก้ไข**

- เก็บ Token ดิบไว้ใน Context ด้วย เพราะการยกเลิกต้องรู้ว่าจะยกเลิกใบไหน

```go
const (
	actorKey = "actor"
	tokenKey = "token"
)

	c.Locals(actorKey, actor)
	c.Locals(tokenKey, token)

func tokenOf(c *fiber.Ctx) string {
	t, _ := c.Locals(tokenKey).(string)
	return t
}
```

**ไฟล์ที่ 3 — `delivery/http/auth.go` สิ่งที่แก้ไข**

```go
func (h Handler) Logout(c *fiber.Ctx) error {
	if err := h.auth.Logout(tokenOf(c)); err != nil {
		return toHTTP(c, err)
	}
	return c.SendStatus(204)
}
```

**ไฟล์ที่ 4 — `delivery/http/error.go` สิ่งที่แก้ไข**

- เพิ่มอีกหนึ่ง `case` ในไฟล์เดียวของโปรแกรมที่ตัดสินเรื่อง Status Code

```go
case errors.Is(err, domain.ErrRevokeUnsupported):
	// 501 แปลว่า "เซิร์ฟเวอร์ยังไม่รองรับสิ่งที่ร้องขอ"
	// ไม่ใช่ความผิดของผู้เรียก แต่เป็นข้อจำกัดของ Adapter ที่ประกอบไว้
	return c.Status(501).JSON(msg(err.Error()))
```

- 📄 ดู `usecase/auth.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package usecase
    
    import (
    	"errors"
    	"fmt"
    	"strings"
    	"sync/atomic"
    
    	"food-api/domain"
    )
    
    // Interface สามตัวนี้ประกาศไว้ที่ชั้น Use Case ซึ่งเป็นผู้ใช้งาน
    // ไม่ใช่ที่ชั้นนอกซึ่งเป็นผู้ให้บริการ หลักการเดียวกับ OrderRepository ในบทที่แล้ว
    //
    // ผลคือชั้นนี้ไม่รู้จัก bcrypt และไม่รู้จัก JWT เลยแม้แต่ชื่อ
    
    type UserRepository interface {
    	CreateUser(u domain.User) error
    	FindUserByEmail(email string) (domain.User, error)
    }
    
    // PasswordHasher ซ่อนวิธีเข้ารหัสไว้ข้างหลัง จะเป็น bcrypt หรือ argon2 ก็ได้
    type PasswordHasher interface {
    	Hash(plain string) (string, error)
    	Compare(hash, plain string) error
    }
    
    // TokenIssuer ซ่อนรูปแบบของ Token ไว้ข้างหลัง จะเป็น JWT หรือแบบอื่นก็ได้
    type TokenIssuer interface {
    	Issue(u domain.User) (string, error)
    	Parse(token string) (domain.User, error)
    }
    
    // TokenRevoker เป็นความสามารถเสริม ไม่ใช่ TokenIssuer ทุกตัวที่ทำได้
    // จึงแยกออกมาเป็น Interface ต่างหาก แล้วค่อยถามตอนใช้งานว่าตัวจริงทำได้หรือไม่
    // รูปแบบนี้พบได้ทั่วไปในไลบรารีมาตรฐานของ Go เช่น io.ReaderFrom
    type TokenRevoker interface {
    	Revoke(token string) error
    }
    
    type AuthUseCase struct {
    	repo   UserRepository
    	hasher PasswordHasher
    	tokens TokenIssuer
    	seq    *int64
    
    	// dummyHash คือ Hash ของรหัสผ่านที่ไม่มีใครใช้จริง
    	// มีไว้เผาเวลาให้เท่ากัน ตอนที่ค้นหาอีเมลแล้วไม่พบ ดูคำอธิบายใน Login
    	dummyHash string
    }
    
    func NewAuthUseCase(r UserRepository, h PasswordHasher, t TokenIssuer) AuthUseCase {
    	var n int64
    	dummy, _ := h.Hash("a-password-that-nobody-will-ever-use")
    	return AuthUseCase{repo: r, hasher: h, tokens: t, seq: &n, dummyHash: dummy}
    }
    
    func (a AuthUseCase) nextID() string {
    	return fmt.Sprintf("u%d", atomic.AddInt64(a.seq, 1))
    }
    
    // Register คือทางสมัครสาธารณะ ผลลัพธ์เป็นลูกค้าเสมอ
    //
    // สังเกตว่าไม่มีพารามิเตอร์ role ให้ส่งเข้ามา ผู้สมัครจึงเลือกสิทธิ์ของตนเองไม่ได้
    // แม้ชั้นนอกจะเผลอส่งค่ามาก็ตาม นี่คือการปิดช่องโหว่ Mass Assignment ที่ต้นทาง
    func (a AuthUseCase) Register(email, password string) (domain.User, error) {
    	return a.create(email, password, domain.Customer)
    }
    
    // CreateStaff ไม่มี Endpoint สาธารณะเส้นใดชี้มาที่นี่
    // ใช้ตอนตั้งระบบเท่านั้น บัญชีพนักงานคนแรกจึงมาจากผู้ดูแลระบบ ไม่ใช่จากการสมัครเอง
    func (a AuthUseCase) CreateStaff(email, password string) (domain.User, error) {
    	return a.create(email, password, domain.Staff)
    }
    
    func (a AuthUseCase) create(email, password string, role domain.Role) (domain.User, error) {
    	email = strings.ToLower(strings.TrimSpace(email))
    
    	if email == "" || !strings.Contains(email, "@") {
    		return domain.User{}, domain.ErrInvalidCredentials
    	}
    	// ความยาวรหัสผ่านเป็นกฎของระบบ จึงตรวจที่ชั้นนี้ ไม่ใช่ที่ Handler
    	if len(password) < 8 {
    		return domain.User{}, domain.ErrWeakPassword
    	}
    	if !role.Valid() {
    		role = domain.Customer
    	}
    
    	if _, err := a.repo.FindUserByEmail(email); err == nil {
    		return domain.User{}, domain.ErrEmailTaken
    	}
    
    	hash, err := a.hasher.Hash(password)
    	if err != nil {
    		return domain.User{}, fmt.Errorf("register: %w", err)
    	}
    
    	// ตัวนับ ID เริ่มจากศูนย์ใหม่ทุกครั้งที่โปรแกรมรีสตาร์ท แต่ข้อมูลในฐานข้อมูลไม่ได้หายไปด้วย
    	// หมายเลขที่ได้จึงอาจชนกับผู้ใช้ที่มีอยู่แล้ว วิธีแก้คือขยับไปหมายเลขถัดไปจนกว่าจะว่าง
    	for i := 0; i < 100; i++ {
    		u := domain.User{ID: a.nextID(), Email: email, PasswordHash: hash, Role: role}
    
    		err := a.repo.CreateUser(u)
    		if err == nil {
    			return u, nil
    		}
    		// ชนเฉพาะเรื่อง ID เท่านั้นที่ลองใหม่ได้ ปัญหาอื่นให้คืนออกไปทันที
    		if !errors.Is(err, domain.ErrIDTaken) {
    			return domain.User{}, fmt.Errorf("register: %w", err)
    		}
    	}
    	return domain.User{}, fmt.Errorf("register: หาหมายเลขผู้ใช้ที่ว่างไม่ได้")
    }
    
    func (a AuthUseCase) Login(email, password string) (string, error) {
    	email = strings.ToLower(strings.TrimSpace(email))
    
    	u, err := a.repo.FindUserByEmail(email)
    	if err != nil {
    		// ตอบข้อความเดียวกันทั้งกรณีไม่มีอีเมลนี้ และกรณีรหัสผ่านผิด
    		// มิฉะนั้นผู้โจมตีจะใช้ API นี้ไล่เดาว่าอีเมลใดมีอยู่ในระบบบ้าง
    		//
    		// ข้อความเหมือนกันอย่างเดียวยังไม่พอ เพราะ bcrypt ใช้เวลาราว 50 มิลลิวินาที
    		// ถ้าเส้นทางนี้คืนค่าทันทีโดยไม่ตรวจอะไรเลย ผู้โจมตีจับเวลาก็แยกออกทันที
    		// จึงต้องเสียเวลาให้เท่ากัน ด้วยการตรวจกับ Hash ปลอมที่ไม่มีวันตรงกับอะไร
    		_ = a.hasher.Compare(a.dummyHash, password)
    		return "", domain.ErrInvalidCredentials
    	}
    
    	if err := a.hasher.Compare(u.PasswordHash, password); err != nil {
    		return "", domain.ErrInvalidCredentials
    	}
    
    	return a.tokens.Issue(u)
    }
    
    // Identify แปลง Token กลับเป็นตัวตน ใช้โดย Middleware ที่ชั้นนอก
    func (a AuthUseCase) Identify(token string) (domain.User, error) {
    	u, err := a.tokens.Parse(token)
    	if err != nil {
    		return domain.User{}, domain.ErrUnauthorized
    	}
    	return u, nil
    }
    
    // Logout ยกเลิก Token ที่ถืออยู่ให้ใช้ไม่ได้ทันที
    //
    // ทำได้จริงเฉพาะเมื่อ Adapter ที่ประกอบเข้ามารองรับการยกเลิก
    // ถ้าเป็น JWT ซึ่งเซิร์ฟเวอร์ไม่ได้จำอะไรไว้เลย จะยกเลิกก่อนหมดอายุไม่ได้
    func (a AuthUseCase) Logout(token string) error {
    	r, ok := a.tokens.(TokenRevoker)
    	if !ok {
    		return domain.ErrRevokeUnsupported
    	}
    	return r.Revoke(token)
    }
    ```
    
- 📄 ดู `delivery/http/middleware_auth.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"strings"
    
    	"food-api/domain"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // actorKey คือชื่อช่องที่ใช้ฝากตัวตนของผู้เรียกไว้ใน Context ของ Request
    // tokenKey เก็บ Token ดิบไว้ด้วย เพราะการออกจากระบบต้องรู้ว่าจะยกเลิกใบไหน
    const (
    	actorKey = "actor"
    	tokenKey = "token"
    )
    
    // RequireAuth เป็น Middleware ที่ทำหน้าที่เดียว คือแปลง Token เป็นตัวตน
    //
    // สังเกตว่าโค้ดที่รู้จักคำว่า Bearer มีอยู่ที่นี่ที่เดียวในโปรเจกต์
    // ชั้น Use Case และชั้น Domain ไม่เคยเห็นคำนี้เลย
    func (h Handler) RequireAuth(c *fiber.Ctx) error {
    	raw := c.Get("Authorization")
    
    	const prefix = "Bearer "
    	if !strings.HasPrefix(raw, prefix) {
    		return toHTTP(c, domain.ErrUnauthorized)
    	}
    
    	token := strings.TrimSpace(strings.TrimPrefix(raw, prefix))
    	if token == "" {
    		return toHTTP(c, domain.ErrUnauthorized)
    	}
    
    	actor, err := h.auth.Identify(token)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    
    	// ฝากตัวตนไว้ให้ Handler ตัวถัดไปหยิบใช้
    	c.Locals(actorKey, actor)
    	c.Locals(tokenKey, token)
    	return c.Next()
    }
    
    // RequireStaff ใช้ต่อจาก RequireAuth สำหรับเส้นทางที่เปิดให้เฉพาะพนักงาน
    func (h Handler) RequireStaff(c *fiber.Ctx) error {
    	if !actorOf(c).IsStaff() {
    		return toHTTP(c, domain.ErrForbidden)
    	}
    	return c.Next()
    }
    
    // actorOf หยิบตัวตนออกจาก Context
    // หากเส้นทางใดไม่ได้ผ่าน RequireAuth มาก่อน จะได้ User ว่างซึ่งไม่มีสิทธิ์ใด
    func actorOf(c *fiber.Ctx) domain.User {
    	u, _ := c.Locals(actorKey).(domain.User)
    	return u
    }
    
    // tokenOf หยิบ Token ดิบออกจาก Context ใช้เฉพาะตอนออกจากระบบ
    func tokenOf(c *fiber.Ctx) string {
    	t, _ := c.Locals(tokenKey).(string)
    	return t
    }
    ```
    
- 📄 ดู `delivery/http/auth.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package http
    
    import (
    	"food-api/domain"
    
    	"github.com/gofiber/fiber/v2"
    )
    
    // DTO ของการสมัครและเข้าสู่ระบบ อยู่ชั้นนอกเช่นเดิม
    // registerRequest ตั้งใจไม่มีฟีลด์ Role
    //
    // ถ้ามีฟีลด์นี้ ผู้ใช้จะส่ง {"role":"staff"} เข้ามาแล้วตั้งตนเป็นพนักงานได้ทันที
    // ช่องโหว่แบบนี้เรียกว่า Mass Assignment คือปล่อยให้ข้อมูลจากภายนอก
    // ไหลเข้าไปตั้งค่าฟีลด์ที่ผู้ใช้ไม่ควรมีสิทธิ์กำหนด
    type registerRequest struct {
    	Email    string `json:"email"`
    	Password string `json:"password"`
    }
    
    type loginRequest struct {
    	Email    string `json:"email"`
    	Password string `json:"password"`
    }
    
    // userResponse ตั้งใจไม่มีฟีลด์ PasswordHash
    // การแปลงผ่าน DTO ทำให้ค่าที่ไม่ควรออกไป ไม่มีทางหลุดออกไปโดยบังเอิญ
    type userResponse struct {
    	ID    string `json:"id"`
    	Email string `json:"email"`
    	Role  string `json:"role"`
    }
    
    type tokenResponse struct {
    	Token string `json:"token"`
    }
    
    func toUserResponse(u domain.User) userResponse {
    	return userResponse{ID: u.ID, Email: u.Email, Role: string(u.Role)}
    }
    
    func (h Handler) Register(c *fiber.Ctx) error {
    	var req registerRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	u, err := h.auth.Register(req.Email, req.Password)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.Status(201).JSON(toUserResponse(u))
    }
    
    func (h Handler) Login(c *fiber.Ctx) error {
    	var req loginRequest
    	if err := c.BodyParser(&req); err != nil {
    		return c.Status(400).JSON(msg("invalid body"))
    	}
    
    	token, err := h.auth.Login(req.Email, req.Password)
    	if err != nil {
    		return toHTTP(c, err)
    	}
    	return c.JSON(tokenResponse{Token: token})
    }
    
    // Me ใช้ตรวจว่า Token ที่ถืออยู่แทนตัวตนใด
    func (h Handler) Me(c *fiber.Ctx) error {
    	return c.JSON(toUserResponse(actorOf(c)))
    }
    
    // Logout ยกเลิก Token ใบที่กำลังใช้อยู่
    //
    // ผลลัพธ์ต่างกันตามโหมดที่ประกอบไว้ใน main.go
    // โหมด session ตอบ 204 แล้ว Token ใช้ต่อไม่ได้ทันที
    // โหมด jwt ตอบ 501 เพราะเซิร์ฟเวอร์ไม่ได้จำอะไรไว้ จึงไม่มีอะไรให้ลบ
    func (h Handler) Logout(c *fiber.Ctx) error {
    	if err := h.auth.Logout(tokenOf(c)); err != nil {
    		return toHTTP(c, err)
    	}
    	return c.SendStatus(204)
    }
    ```
    
- 📄 ดู `delivery/http/error.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
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
    
    	// เพิ่มในบทนี้ ยังคงเป็นไฟล์เดียวของทั้งโปรแกรมที่ตัดสินเรื่อง Status Code
    	case errors.Is(err, domain.ErrUnauthorized), errors.Is(err, domain.ErrInvalidCredentials):
    		// 401 แปลว่า "ยังไม่รู้ว่าคุณเป็นใคร"
    		c.Set("WWW-Authenticate", "Bearer")
    		return c.Status(401).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrForbidden):
    		// 403 แปลว่า "รู้แล้วว่าคุณเป็นใคร แต่คุณไม่มีสิทธิ์"
    		return c.Status(403).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrEmailTaken):
    		return c.Status(409).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrWeakPassword):
    		return c.Status(422).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrRevokeUnsupported):
    		// 501 แปลว่า "เซิร์ฟเวอร์ยังไม่รองรับสิ่งที่ร้องขอ"
    		// ไม่ใช่ความผิดของผู้เรียก แต่เป็นข้อจำกัดของ Adapter ที่ประกอบไว้
    		return c.Status(501).JSON(msg(err.Error()))
    
    	case errors.Is(err, domain.ErrTooManyRequests):
    		return c.Status(429).JSON(msg(err.Error()))
    
    	default:
    		// Error ที่ไม่รู้จัก หมายความว่าเป็นข้อผิดพลาดของระบบเราเอง
    		// บันทึกไว้ตรวจสอบภายหลัง แต่ไม่ส่งให้ผู้ใช้ เพราะอาจมีชื่อตารางปะปนอยู่
    		log.Println(err)
    		return c.Status(500).JSON(msg("internal error"))
    	}
    }
    ```
    

**ผลลัพธ์**

```bash
go build ./...
```

<aside>
⚠️

**นี่คือการแก้ชั้นในเพียงจุดเดียวของทั้งหัวข้อนี้ และต้องเข้าใจให้ตรงว่าแก้เพราะอะไร**

ไม่ได้แก้เพราะจะเพิ่ม Session แต่แก้เพราะจะเพิ่ม **ความสามารถใหม่** ชื่อการออกจากระบบ ซึ่งเป็นคนละเรื่องกัน

การสลับ JWT เป็น Session ในขั้นถัดไป จะไม่แตะไฟล์เหล่านี้อีกเลยแม้แต่บรรทัดเดียว

</aside>

#### ขั้นที่ 17 — สลับทั้งระบบด้วยการแก้ .env บรรทัดเดียว

<aside>
✏️

**ขั้นนี้เป็นการแก้ไฟล์เดิม ไม่ใช่การสร้างไฟล์ใหม่**

**1. เปิดไฟล์** `main.go` แก้ตามรายการด้านล่าง หรือลบทั้งไฟล์แล้ววางฉบับเต็มทับก็ได้

**2. เปิดไฟล์** `.env.example` เพิ่มสามบรรทัดท้ายไฟล์ แล้วคัดลอกลง `.env` ของตนเอง

**3. รัน** `go run .`

คำสั่งทั้งหมดในขั้นนี้ ให้รันในโฟลเดอร์ `food-api`

</aside>

**สิ่งที่แก้ไขใน `main.go`**

- เพิ่มตัวแปร `sessions` ไว้ข้าง ๆ `orders` และ `users` แล้วเลือก Adapter ในที่เดียวกัน
- เปลี่ยนบรรทัดที่สร้าง `tokens` ให้เลือกตามค่า `AUTH_MODE`
- เพิ่มเส้นทาง `POST /api/auth/logout` ไว้ในกลุ่มที่ต้องเข้าสู่ระบบก่อน

```go
var tokens usecase.TokenIssuer
switch mode := envOr("AUTH_MODE", "jwt"); mode {
case "session":
	log.Println("ใช้ Session แบบ Stateful ยกเลิกได้ทันที แต่ต้องอ่านที่เก็บข้อมูลทุก Request")
	tokens = security.NewSessionIssuer(sessions, ttl)
case "jwt":
	log.Println("ใช้ JWT แบบ Stateless ไม่ต้องอ่านที่เก็บข้อมูล แต่ยกเลิกก่อนหมดอายุไม่ได้")
	tokens = security.NewJWTIssuer(secret, ttl)
default:
	log.Fatal("AUTH_MODE ต้องเป็น jwt หรือ session เท่านั้น ได้รับ: ", mode)
}
```

**สิ่งที่เพิ่มใน `.env.example` และ `.env`**

```bash
AUTH_MODE=jwt
```

ขั้นนี้มีการแก้ไขหลายจุดในไฟล์เดียว วิธีที่ปลอดภัยที่สุดคือ **ลบเนื้อหาเดิมทั้งหมด แล้ววางไฟล์ฉบับเต็มด้านล่างทับลงไป**

- 📄 ดู `main.go` ฉบับเต็มหลังแก้ไข — ใช้วางทับทั้งไฟล์ได้
    
    ```go
    package main
    
    import (
    	"context"
    	"log"
    	"os"
    	"strings"
    	"time"
    
    	"food-api/delivery/http"
    	"food-api/domain"
    	"food-api/repository/memory"
    	"food-api/repository/postgres"
    	"food-api/security"
    	"food-api/usecase"
    
    	"github.com/gofiber/fiber/v2"
    	"github.com/gofiber/fiber/v2/middleware/cors"
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
    
    // เป็นไฟล์เดียวในโปรเจกต์ที่อ้างถึงทุกชั้นพร้อมกัน และเป็นที่เดียวที่เลือกว่า
    // จะใช้ Adapter ตัวใด ทั้งที่เก็บข้อมูล วิธี Hash รหัสผ่าน และรูปแบบ Token
    func main() {
    	_ = godotenv.Load()
    
    	// ค่าลับต้องมาจาก Environment เสมอ และต้องหยุดทำงานทันทีเมื่อไม่มี
    	// ห้ามใส่ค่าสำรองไว้ในโค้ด เพราะจะกลายเป็นรหัสลับที่ทุกคนบนโลกรู้
    	secret := os.Getenv("JWT_SECRET")
    	if len(strings.TrimSpace(secret)) < 32 {
    		log.Fatal("ต้องกำหนด JWT_SECRET ใน .env และต้องยาวอย่างน้อย 32 ตัวอักษร")
    	}
    
    	ttl, err := time.ParseDuration(envOr("JWT_TTL", "1h"))
    	if err != nil {
    		log.Fatal("JWT_TTL ไม่ถูกต้อง: ", err)
    	}
    
    	// 1. ที่เก็บข้อมูล เลือกจาก Environment เช่นเดิม
    	var orders usecase.OrderRepository
    	var users usecase.UserRepository
    	var sessions security.SessionStore
    
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
    		orders = postgres.NewOrderRepo(pool)
    		users = postgres.NewUserRepo(pool)
    		sessions = postgres.NewSessionRepo(pool)
    	} else {
    		log.Println("ใช้ Memory")
    		orders = memory.NewOrderRepo(seedMenu())
    		users = memory.NewUserRepo()
    		sessions = memory.NewSessionRepo()
    	}
    
    	// 2. Adapter ด้านความปลอดภัย
    	//    ชั้น Use Case รู้จักเพียง Interface ไม่รู้จักชื่อ bcrypt และ JWT
    	hasher := security.NewBcryptHasher(0)
    
    	// นี่คือจุดสลับทั้งหมดของระบบยืนยันตัวตน
    	//
    	// ทั้งสองตัวทำตาม Interface usecase.TokenIssuer เหมือนกัน
    	// การเปลี่ยนค่า AUTH_MODE ใน .env จึงเปลี่ยนวิธียืนยันตัวตนทั้งระบบ
    	// โดยที่ชั้น Domain ชั้น Use Case และ Handler ไม่ถูกแตะแม้แต่บรรทัดเดียว
    	var tokens usecase.TokenIssuer
    	switch mode := envOr("AUTH_MODE", "jwt"); mode {
    	case "session":
    		log.Println("ใช้ Session แบบ Stateful ยกเลิกได้ทันที แต่ต้องอ่านที่เก็บข้อมูลทุก Request")
    		tokens = security.NewSessionIssuer(sessions, ttl)
    	case "jwt":
    		log.Println("ใช้ JWT แบบ Stateless ไม่ต้องอ่านที่เก็บข้อมูล แต่ยกเลิกก่อนหมดอายุไม่ได้")
    		tokens = security.NewJWTIssuer(secret, ttl)
    	default:
    		log.Fatal("AUTH_MODE ต้องเป็น jwt หรือ session เท่านั้น ได้รับ: ", mode)
    	}
    
    	// 3. กฎ
    	uc := usecase.NewOrderUseCase(orders)
    	auth := usecase.NewAuthUseCase(users, hasher, tokens)
    
    	// บัญชีพนักงานคนแรกมาจากผู้ดูแลระบบ ไม่ใช่จากการสมัครเอง
    	// เส้นทาง /auth/register สาธารณะสร้างได้เฉพาะลูกค้าเท่านั้น
    	if email, pw := os.Getenv("STAFF_EMAIL"), os.Getenv("STAFF_PASSWORD"); email != "" && pw != "" {
    		if _, err := auth.CreateStaff(email, pw); err != nil {
    			log.Println("ข้ามการสร้างบัญชีพนักงานเริ่มต้น:", err)
    		} else {
    			log.Println("สร้างบัญชีพนักงานเริ่มต้นแล้ว:", email)
    		}
    	}
    
    	// 4. ชั้นขนส่ง
    	h := http.NewHandler(uc, auth)
    
    	app := fiber.New(fiber.Config{
    		// จำกัดขนาด Body ป้องกันการส่งไฟล์ใหญ่มากเข้ามาจนหน่วยความจำหมด
    		BodyLimit: 1 * 1024 * 1024,
    		// ตัดการเชื่อมต่อที่ค้างนานเกินไป
    		ReadTimeout:  10 * time.Second,
    		WriteTimeout: 10 * time.Second,
    	})
    
    	app.Use(http.SecurityHeaders)
    
    	// อนุญาตเฉพาะต้นทางที่ระบุไว้ ห้ามใช้ * คู่กับการส่ง Credential
    	app.Use(cors.New(cors.Config{
    		AllowOrigins: envOr("CORS_ORIGINS", "http://localhost:5173"),
    		AllowMethods: "GET,POST,PATCH,DELETE",
    		AllowHeaders: "Content-Type,Authorization",
    	}))
    
    	app.Use(func(c *fiber.Ctx) error {
    		err := c.Next()
    		log.Println(c.Method(), c.Path(), c.Response().StatusCode())
    		return err
    	})
    
    	// ตัวนับสองชุด ชุดกว้างสำหรับทุกเส้นทาง และชุดเข้มเฉพาะหน้าเข้าสู่ระบบ
    	general := http.NewRateLimiter(100, time.Minute)
    	loginer := http.NewRateLimiter(5, time.Minute)
    	app.Use(general.Limit)
    
    	api := app.Group("/api")
    
    	// ── เปิดให้ทุกคน ──
    	api.Post("/auth/register", h.Register)
    	api.Post("/auth/login", loginer.Limit, h.Login)
    	api.Get("/menu", h.ListMenu)
    
    	// ── ต้องเข้าสู่ระบบก่อน ──
    	me := api.Group("", h.RequireAuth)
    	me.Get("/auth/me", h.Me)
    	me.Post("/auth/logout", h.Logout)
    	me.Post("/orders", h.PlaceOrder)
    	me.Get("/orders/:id", h.GetOrder)
    	me.Get("/orders/:id/total", h.OrderTotal)
    	me.Post("/orders/:id/cancel", h.CancelOrder)
    
    	// ── เฉพาะพนักงานร้าน ──
    	staff := api.Group("", h.RequireAuth, h.RequireStaff)
    	staff.Patch("/orders/:id/status", h.SetStatus)
    
    	log.Fatal(app.Listen(":3000"))
    }
    
    func envOr(key, fallback string) string {
    	if v := os.Getenv(key); v != "" {
    		return v
    	}
    	return fallback
    }
    ```
    

**ผลลัพธ์**

```bash
go run .
```

```bash
ใช้ Postgres
ใช้ JWT แบบ Stateless ไม่ต้องอ่านที่เก็บข้อมูล แต่ยกเลิกก่อนหมดอายุไม่ได้
สร้างบัญชีพนักงานเริ่มต้นแล้ว: staff@example.com
```

**ทดลองออกจากระบบในโหมด `jwt` ซึ่งเป็นค่าตั้งต้น**

```bash
TOKEN=$(curl -s -X POST localhost:3000/api/auth/login   -H 'Content-Type: application/json'   -d '{"email":"alice@example.com","password":"password123"}' | jq -r .token)
curl -s -i -X POST localhost:3000/api/auth/logout -H "Authorization: Bearer $TOKEN" | head -1
```

```bash
HTTP/1.1 501 Not Implemented
```

Token ใบเดิมยังใช้งานได้อยู่ เพราะเซิร์ฟเวอร์ไม่ได้จำอะไรไว้ จึงไม่มีอะไรให้ลบ

```bash
curl -s -o /dev/null -w "%{http_code}\n" localhost:3000/api/auth/me -H "Authorization: Bearer $TOKEN"
```

```bash
200
```

**คราวนี้แก้ `.env` เป็น `AUTH_MODE=session` แล้วรีสตาร์ทเซิร์ฟเวอร์ ทำซ้ำทั้งชุด**

```bash
ใช้ Session แบบ Stateful ยกเลิกได้ทันที แต่ต้องอ่านที่เก็บข้อมูลทุก Request
```

สิ่งแรกที่สังเกตได้คือหน้าตาของ Token เปลี่ยนไปคนละแบบ

```bash
echo $TOKEN
```

```bash
vWn83HWc8zXh_YUnExldMsZvg_gNE4QNm3XEVMNzpZ4
```

สั้นลงจาก 204 ตัวอักษรเหลือ 43 ตัวอักษร และอ่านอะไรไม่ได้เลย เพราะเป็นเพียงกุญแจไว้เปิดหาข้อมูลจริงในตาราง

```bash
docker compose exec db psql -U app -d food -c "SELECT id, user_id, email, role FROM sessions;"
```

```bash
                     id                      | user_id |       email       |   role
---------------------------------------------+---------+-------------------+----------
 vWn83HWc8zXh_YUnExldMsZvg_gNE4QNm3XEVMNzpZ4 | u2      | alice@example.com | customer
(1 row)
```

ทีนี้ลองออกจากระบบอีกครั้ง

```bash
curl -s -i -X POST localhost:3000/api/auth/logout -H "Authorization: Bearer $TOKEN" | head -1
```

```bash
HTTP/1.1 204 No Content
```

แล้วลองใช้ Token ใบเดิมซ้ำ

```bash
curl -s -w "\nHTTP %{http_code}\n" localhost:3000/api/auth/me -H "Authorization: Bearer $TOKEN"
```

```bash
{"error":"unauthorized"}
HTTP 401
```

<aside>
✅

**นี่คือผลลัพธ์ทั้งหมดของหัวข้อนี้ในภาพเดียว**

Endpoint เดิม โค้ดชั้น Domain เดิม ชั้น Use Case เดิม และ Handler เดิมทุกบรรทัด

เปลี่ยนเพียงค่าใน `.env` หนึ่งบรรทัด พฤติกรรมด้านความปลอดภัยของทั้งระบบก็เปลี่ยนไปคนละแบบ

ลองค้นคำว่า `Session` ในโฟลเดอร์ `domain`, `usecase` และ `delivery` จะไม่พบเลยแม้แต่ที่เดียว

นั่นคือความหมายของประโยคที่พูดไว้ตั้งแต่ต้นวันว่า **JWT เป็นเพียงรายละเอียดของชั้นนอก**

</aside>

### 3.11.1 เทียบข้อดีข้อเสีย จากผลที่เพิ่งเห็นกับตา

สองขั้นที่ผ่านมาให้หลักฐานครบแล้ว นำมาเรียงเป็นตารางได้ดังนี้

| ประเด็น | JWT | Session |
| --- | --- | --- |
| ความยาว Token ที่วัดได้จริง | 204 ตัวอักษร และยาวขึ้นเมื่อใส่ข้อมูลเพิ่ม | 43 ตัวอักษร คงที่เสมอ |
| อ่านเนื้อหาข้างในได้ไหม | ได้ ใครก็ถอด Base64 อ่านได้ **ห้ามใส่ความลับ** | ไม่ได้ เป็นค่าสุ่มล้วน ไม่มีความหมายนอกเซิร์ฟเวอร์ |
| ตรวจหนึ่ง Request ใช้อะไร | คำนวณลายเซ็น ไม่แตะฐานข้อมูล | อ่านฐานข้อมูล 1 ครั้ง |
| ออกจากระบบ | **ทำไม่ได้** ตอบ 501 และ Token ยังใช้ได้ต่อ | ตอบ 204 แล้ว Token ตายทันที กลายเป็น 401 |
| เพิ่มสำเนาเซิร์ฟเวอร์ | ทุกสำเนาตรวจได้เลย ไม่ต้องคุยกัน | ทุกสำเนาต้องต่อที่เก็บกลางตัวเดียวกัน |
| ฐานข้อมูลล่ม | ผู้ที่ถือ Token อยู่แล้วยังเรียกได้ | ทุกคนใช้งานไม่ได้ทันที |
| เปลี่ยนรหัสผ่านแล้วเตะทุกเครื่อง | ทำไม่ได้ ต้องรอ Token หมดอายุ | ลบทุกแถวของผู้ใช้คนนั้น จบทันที |
| Token หลุดไปอยู่กับผู้โจมตี | ใช้ได้จนหมดอายุ หยุดไม่ได้ | ลบแถวเดียว ตายทันที |
| ต้นทุนที่เพิ่มขึ้น | ไม่มีตารางเพิ่ม | มีตาราง `sessions` ที่ต้องดูแลและกวาดล้าง |

<aside>
💡

อ่านตารางแล้วจะเห็นว่า **ข้อดีของฝั่งหนึ่งคือข้อเสียของอีกฝั่งเสมอ** ไม่มีตัวใดชนะขาด

แก่นของเรื่องมีประโยคเดียวคือ *ใครเป็นคนจำว่าใครล็อกอินอยู่*

ถ้าเซิร์ฟเวอร์ไม่จำ ก็ขยายง่ายแต่ยกเลิกไม่ได้  ถ้าเซิร์ฟเวอร์จำ ก็ยกเลิกได้แต่ต้องมีที่เก็บกลาง

</aside>

### 3.11.2 แล้วงานจริงควรเลือกอันไหน

คำถามที่ควรถามไม่ใช่ *อันไหนดีกว่า* แต่เป็น *ระบบนี้ยอมให้ยกเลิกช้าได้แค่ไหน*

| เลือก **JWT** เมื่อ | เหตุผล |
| --- | --- |
| ผู้เรียกเป็นโปรแกรม ไม่ใช่คน เช่น Service เรียกหากัน | ไม่มีปุ่มออกจากระบบให้กดอยู่แล้ว |
| ระบบขยายเป็นหลายสำเนา และไม่อยากมีที่เก็บกลาง | ทุกสำเนาตรวจ Token ได้เองโดยไม่ต้องคุยกัน |
| API สาธารณะที่ให้บุคคลที่สามเรียก | ผู้เรียกไม่ได้อยู่ในระบบเรา จึงเก็บ Session ให้ไม่ได้ |
| ยอมรับได้ว่าการยกเลิกช้าเท่าอายุ Token | จึงต้องตั้งอายุสั้น ระดับ 5 ถึง 15 นาที ไม่ใช่ 1 ชั่วโมงอย่างที่ตั้งไว้ในบทนี้ |

| เลือก **Session** เมื่อ | เหตุผล |
| --- | --- |
| เป็นเว็บที่คนใช้ และมีปุ่มออกจากระบบ | ผู้ใช้คาดหวังว่ากดแล้วต้องออกจริง ไม่ใช่ออกในอีกหนึ่งชั่วโมง |
| ต้องมีปุ่ม ออกจากระบบทุกอุปกรณ์ | ลบทุกแถวของผู้ใช้คนนั้นได้ในคำสั่งเดียว |
| ต้องแบนผู้ใช้ หรือตัดสิทธิ์ได้ทันที | ธนาคาร ระบบสุขภาพ และงานที่มีข้อกำหนดกำกับดูแล มักบังคับข้อนี้ |
| เปลี่ยนรหัสผ่านแล้วต้องเตะเครื่องอื่นออก | เป็นพฤติกรรมมาตรฐานที่ผู้ใช้คาดหวัง |

<aside>
✅

**ระบบจริงขนาดใหญ่มักไม่เลือกข้างใดข้างหนึ่ง แต่ใช้ทั้งคู่**

**Access Token** เป็น JWT อายุสั้นมาก ระดับ 5 ถึง 15 นาที ใช้ยิง API ทุกครั้งโดยไม่ต้องแตะฐานข้อมูล

**Refresh Token** เป็น Session ที่เก็บในฐานข้อมูล อายุยาว ใช้ขอ Access Token ใบใหม่ และ **ยกเลิกได้ทันที**

ผลคือได้ความเร็วของ JWT ในเส้นทางที่เรียกบ่อย และได้อำนาจยกเลิกของ Session ในเส้นทางที่เรียกนาน ๆ ครั้ง

หน้าต่างความเสี่ยงเมื่อ Token หลุด จึงเหลือเท่ากับอายุของ Access Token เท่านั้น ไม่ใช่ทั้งชั่วโมง

โครงสร้างที่เขียนไว้ในบทนี้รองรับการทำแบบนั้นได้ทันที เพราะ `TokenIssuer` ตัวที่สามจะเสียบเข้ามาโดยไม่กระทบชั้นใดเลย

</aside>

### 3.11.3 อีกคำถามที่ต้องตอบคู่กันเสมอ คือฝั่งหน้าเว็บเก็บ Token ไว้ที่ไหน

เลือกวิธียืนยันตัวตนแล้วยังไม่จบ เพราะที่เก็บ Token ฝั่งเบราว์เซอร์มีผลต่อความปลอดภัยไม่แพ้กัน และเป็นคนละเรื่องกับที่เพิ่งคุยไป

|  | เก็บใน `localStorage` | เก็บใน Cookie แบบ `HttpOnly` |
| --- | --- | --- |
| JavaScript อ่านได้ไหม | อ่านได้ จึงเป็นเป้าของ XSS โดยตรง | อ่านไม่ได้ แม้โค้ดที่ถูกฝังมาก็อ่านไม่ได้ |
| ถ้าเว็บโดน XSS | Token ถูกขโมยไปทั้งใบ | ขโมยตัว Token ไม่ได้ แต่ยังสั่งงานแทนผู้ใช้ได้ |
| เสี่ยง CSRF ไหม | ไม่เสี่ยง เพราะต้องแนบ Header เอง | เสี่ยง เพราะเบราว์เซอร์แนบ Cookie ให้อัตโนมัติ |
| สิ่งที่ต้องทำเพิ่ม | ป้องกัน XSS ให้ดี | ตั้ง `SameSite=Lax` หรือ `Strict` และมี CSRF Token |

ระบบที่เราเขียนวันนี้ส่ง Token ผ่าน Header `Authorization` ซึ่งเข้าคู่กับการเก็บใน `localStorage` จึงไม่ต้องกังวลเรื่อง CSRF แต่ต้องระวัง XSS เป็นพิเศษ

<aside>
⚠️

ข้อสรุปที่ใช้ได้จริงคือ **ไม่มีที่เก็บใดปลอดภัยสมบูรณ์** มีแต่การเลือกว่าจะรับความเสี่ยงด้านใด

และไม่ว่าเลือกทางใด สิ่งที่ต้องทำเหมือนกันคือ ตั้งอายุ Token ให้สั้น และเปิด HTTPS เสมอในระบบจริง

</aside>

## 3.12 การโจมตีที่ต้องรู้จัก

ทั้งวันที่ผ่านมาเป็นการสร้าง *เครื่องป้องกัน* หัวข้อนี้กลับด้าน คือดูจากมุมของผู้โจมตีว่าเขาจะลองอะไรบ้าง และโค้ดบรรทัดใดของเราที่หยุดเขาไว้

<aside>
💡

ทุกคำสั่งในหัวข้อนี้ยิงใส่เครื่องของตนเอง บนข้อมูลที่ตนสร้างขึ้นเอง

การเข้าใจว่าการโจมตีทำงานอย่างไร คือเงื่อนไขที่ขาดไม่ได้ของการเขียนโค้ดที่กันมันได้

</aside>

### 3.12.1 ตารางสรุป

| ชื่อการโจมตี | ผู้โจมตีทำอะไร | อะไรหยุดเขาไว้ |
| --- | --- | --- |
| **Broken Object Level Authorization** หรือ IDOR | เปลี่ยนเลขที่ออเดอร์ใน URL เป็นของคนอื่น | `CanBeViewedBy` ใน `domain` พิสูจน์ไว้แล้วในขั้นที่ 13 |
| **Privilege Escalation** ผ่าน Mass Assignment | ส่ง `"role":"staff"` ตอนสมัคร | `Register` ไม่รับพารามิเตอร์ `role` และ DTO ไม่มีฟีลด์นั้น |
| **Token Forgery** | แก้ `role` ใน Payload แล้วส่งกลับมา | การตรวจลายเซ็นใน `security/jwt.go` พิสูจน์ไว้ในขั้นที่ 13 |
| **Algorithm Confusion** หรือ alg none | ส่ง Token ที่บอกว่าไม่ต้องตรวจลายเซ็น | `jwt.WithValidMethods([]string{"HS256"})` |
| **SQL Injection** | แทรกคำสั่ง SQL ลงในค่าที่กรอก | `$1` `$2` ของ pgx ที่ใช้มาตั้งแต่ Day 2 |
| **Brute Force** และ **Credential Stuffing** | ไล่ยิงรหัสผ่านซ้ำ ๆ | `NewRateLimiter(5, time.Minute)` ในขั้นที่ 12 |
| **User Enumeration** | ดูว่าอีเมลใดมีอยู่ในระบบ | ข้อความ 401 เหมือนกัน และการเผาเวลาให้เท่ากัน |
| **Existence Oracle** ผ่าน 403 | แยกว่าออเดอร์มีอยู่จริงไหม จาก Status Code | ตอบ 404 ไม่ใช่ 403 ในขั้นที่ 4 |
| **Clickjacking** และ **MIME Sniffing** | ครอบหน้าเว็บเราไว้ใน iframe | `X-Frame-Options` และ `X-Content-Type-Options` ในขั้นที่ 11 |
| **Denial of Service** จาก Body ขนาดใหญ่ | ส่งไฟล์ขนาดมหาศาลเข้ามา | `BodyLimit` และ Timeout ในขั้นที่ 12 |

### 3.12.2 SQL Injection ลงมือทดลองจริง

หัวข้อนี้แยกออกมาเพราะเป็นการโจมตีที่เก่าที่สุด พบบ่อยที่สุด และเสียหายหนักที่สุด

**ต้นเหตุคือการเอาค่าจากผู้ใช้ไปต่อสายอักขระเข้ากับคำสั่ง SQL** หน้าตาของโค้ดที่ผิดเป็นแบบนี้

```go
// ห้ามเขียนแบบนี้เด็ดขาด
id := c.Params("id")
query := "SELECT id, status, total FROM orders WHERE id = '" + id + "'"
rows, _ := pool.Query(context.Background(), query)
```

ถ้าผู้ใช้ส่ง `id` เป็น `o1'; DROP TABLE orders; --` สิ่งที่ฐานข้อมูลได้รับจะกลายเป็นสองคำสั่ง

```sql
SELECT id, status, total FROM orders WHERE id = 'o1'; DROP TABLE orders; --'
```

ตารางหายไปทั้งตาราง เพราะฐานข้อมูลไม่มีทางรู้ว่าส่วนใดเป็นคำสั่งที่โปรแกรมเมอร์ตั้งใจ และส่วนใดเป็นค่าที่ผู้ใช้กรอก

**โค้ดที่เราเขียนมาตั้งแต่ Day 2 ไม่มีปัญหานี้ เพราะใช้ `$1` มาตลอด**

```go
err := r.db.QueryRow(context.Background(),
	`SELECT id, email, password_hash, role FROM users WHERE email = $1`, email,
).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
```

<aside>
💡

**เหตุผลที่ `$1` ปลอดภัย ไม่ใช่เพราะมันกรองอักขระอันตรายทิ้ง**

แต่เพราะคำสั่งกับค่าถูกส่งไปยังฐานข้อมูล **คนละทาง**

ฐานข้อมูลแปลคำสั่งจนเสร็จก่อน แล้วค่อยรับค่าเข้ามาเติมในช่องที่เตรียมไว้

ตอนที่ค่ามาถึง การแปลคำสั่งจบไปแล้ว ค่าจึงไม่มีทางกลายเป็นคำสั่งได้อีก ต่อให้ข้างในเป็นคำว่า `DROP TABLE` ก็ตาม

</aside>

ทดลองด้วยตนเองสามข้อ ให้เซิร์ฟเวอร์ทำงานค้างไว้ แล้วเรียกทีละคำสั่ง

คำสั่งที่ 1 แทรกคำสั่งลบตารางผ่าน URL

```bash
curl -s -w "\nHTTP %{http_code}\n" \
  "localhost:3000/api/orders/o1%27%3B%20DROP%20TABLE%20orders%3B%20--" \
  -H "Authorization: Bearer $TOKEN"
```

```bash
{"error":"not found"}
HTTP 404
```

ค่าทั้งก้อนถูกใช้เป็น *ชื่อออเดอร์* ตรง ๆ และไม่มีออเดอร์ชื่อนั้น จึงได้ 404 ตามปกติ ตารางยังอยู่ครบ

คำสั่งที่ 2 ลองเงื่อนไขที่เป็นจริงเสมอ ตอนเข้าสู่ระบบ

```bash
curl -s -w "\nHTTP %{http_code}\n" -X POST localhost:3000/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com'"'"' OR '"'"'1'"'"'='"'"'1","password":"wrong"}'
```

```bash
{"error":"invalid email or password"}
HTTP 401
```

ในระบบที่มีช่องโหว่ คำสั่งนี้จะทำให้เข้าสู่ระบบได้โดยไม่ต้องรู้รหัสผ่าน ที่นี่มันเป็นเพียงอีเมลแปลก ๆ ที่หาไม่เจอ

คำสั่งที่ 3 สมัครสมาชิกด้วยชื่อในตำนาน

```bash
curl -s -w "\nHTTP %{http_code}\n" -X POST localhost:3000/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"bobby@tables.com'"'"'; DROP TABLE users; --","password":"password123"}'
```

```bash
{"id":"u3","email":"bobby@tables.com'; drop table users; --","role":"customer"}
HTTP 201
```

เปิดดูในฐานข้อมูลจะพบว่า คำสั่งลบตารางถูกเก็บไว้เป็น **ข้อความธรรมดา** ในช่องอีเมล

```bash
docker compose exec db psql -U app -d food -c "SELECT id, email FROM users WHERE email LIKE 'bobby%';"
```

```bash
 id |                  email
----+-----------------------------------------
 u3 | bobby@tables.com'; drop table users; --
(1 row)
```

<aside>
⚠️

**สองเรื่องที่ต้องระวังต่อ**

**ข้อ 1 ไลบรารีไม่ได้ช่วยอัตโนมัติ** ต่อให้ใช้ ORM หรือ Query Builder ถ้าเผลอใช้ `fmt.Sprintf` ต่อสายอักขระเข้าไปในคำสั่ง ช่องโหว่ก็กลับมาทันที ตัวชี้วัดคือ **ค่าจากผู้ใช้ต้องไม่เคยอยู่ในสายอักขระของคำสั่ง SQL**

**ข้อ 2 บางอย่างใส่เป็นพารามิเตอร์ไม่ได้** เช่นชื่อตาราง ชื่อคอลัมน์ และทิศทางของ `ORDER BY` สิ่งเหล่านี้ `$1` ช่วยไม่ได้

ถ้าจำเป็นต้องให้ผู้ใช้เลือกคอลัมน์ที่ใช้เรียงลำดับ ต้องเทียบกับ **รายชื่อที่อนุญาตไว้ล่วงหน้า** เท่านั้น ห้ามเอาค่าที่ได้รับมาต่อลงไปตรง ๆ

</aside>

### 3.12.3 Mass Assignment ลองตั้งตนเป็นพนักงาน

ลองสมัครสมาชิกพร้อมแอบใส่ `role` เข้าไปด้วย

```bash
curl -s -X POST localhost:3000/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"attacker@example.com","password":"password123","role":"staff"}'
```

```bash
{"id":"u4","email":"attacker@example.com","role":"customer"}
```

ส่ง `staff` เข้าไป แต่ได้ `customer` กลับมา เพราะสองชั้นช่วยกันปิดไว้

- `registerRequest` ไม่มีฟีลด์ `Role` ค่าที่ส่งมาจึงถูกทิ้งตั้งแต่ตอนแปลง JSON
- `Register` ในชั้น Use Case ไม่รับพารามิเตอร์ `role` เลย และเรียก `create` ด้วย `domain.Customer` ตายตัว

<aside>
🚨

**ถ้าเขียนแบบด้านล่างนี้ ระบบจะพัง และเป็นความผิดพลาดที่พบบ่อยมากในงานจริง**

`Role string` json:"role" ใน DTO บวกกับ `Register(req.Email, req.Password, domain.Role(req.Role))`

ผลคือใครก็ได้บนอินเทอร์เน็ตตั้งตนเป็นพนักงานได้ แล้วเปลี่ยนสถานะออเดอร์ของทุกคนได้ทันที

**หลักที่ใช้กันคือ ฟีลด์ที่ผู้ใช้ไม่ควรกำหนดเอง ต้องไม่มีอยู่ใน DTO ตั้งแต่แรก** ไม่ใช่ใส่ไว้แล้วค่อยกรองทีหลัง

บัญชีพนักงานคนแรกจึงมาจาก `STAFF_EMAIL` ใน `.env` ซึ่งมีเฉพาะผู้ดูแลระบบที่แก้ได้

</aside>

### 3.12.4 การโจมตีด้วยการจับเวลา

ขั้นที่ 3 ทำให้ข้อความตอบกลับเหมือนกันทั้งกรณีไม่มีอีเมล และกรณีรหัสผ่านผิด แต่ **ข้อความเหมือนกันอย่างเดียวยังไม่พอ**

เพราะ bcrypt ตั้งใจออกแบบมาให้ช้า ราว 50 มิลลิวินาที ถ้าเส้นทาง *ไม่พบอีเมล* คืนค่าทันทีโดยไม่ตรวจอะไรเลย ผู้โจมตีเพียงจับเวลาก็แยกออกได้ทันที

ตัวเลขที่วัดได้จริงก่อนแก้

```bash
อีเมลที่ไม่มีในระบบ            0.0 ms ต่อครั้ง
อีเมลที่มีจริง รหัสผ่านผิด    48.7 ms ต่อครั้ง
```

ต่างกันเกือบ 50 มิลลิวินาที ซึ่งมากพอที่จะวัดผ่านเครือข่ายได้สบาย ผู้โจมตีจึงไล่ตรวจได้ว่าอีเมลใดเป็นสมาชิกของเราบ้าง

วิธีแก้คือ **เผาเวลาให้เท่ากันทั้งสองทาง** ด้วยการตรวจกับ Hash ปลอมที่ไม่มีวันตรงกับอะไร

```go
u, err := a.repo.FindUserByEmail(email)
if err != nil {
	_ = a.hasher.Compare(a.dummyHash, password)
	return "", domain.ErrInvalidCredentials
}
```

ตัวเลขหลังแก้

```bash
อีเมลที่ไม่มีในระบบ           48.5 ms ต่อครั้ง
อีเมลที่มีจริง รหัสผ่านผิด    49.6 ms ต่อครั้ง
```

<aside>
💡

บทเรียนของหัวข้อนี้ใหญ่กว่าเรื่อง bcrypt

**ระบบรั่วข้อมูลออกไปได้หลายทาง ไม่ใช่แค่ทางเนื้อความที่ตอบกลับ** เวลาที่ใช้ ขนาดของคำตอบ และลำดับการทำงาน ล้วนบอกอะไรบางอย่างได้ทั้งสิ้น

ช่องทางแบบนี้เรียกรวมกันว่า Side Channel และเป็นเหตุผลที่การเปรียบเทียบค่าลับต้องใช้ฟังก์ชันที่ใช้เวลาคงที่เสมอ

</aside>

### 3.12.5 สิ่งที่ยังไม่ได้ทำ และควรทำต่อในระบบจริง

| เรื่อง | ต้องทำอะไรเพิ่ม |
| --- | --- |
| HTTPS | ทั้งบทนี้ทดสอบบน `http` ในระบบจริง Token ที่วิ่งบน `http` ถูกดักอ่านได้ทั้งใบ ต้องบังคับ `https` และเพิ่ม Header `Strict-Transport-Security` |
| Refresh Token | ตั้งอายุ Access Token ให้สั้นระดับนาที แล้วออกใบใหม่ด้วย Refresh Token ที่ยกเลิกได้ ตามที่อธิบายไว้ในหัวข้อ 3.11.2 |
| กวาดล้าง Session ที่หมดอายุ | ตาราง `sessions` โตขึ้นเรื่อย ๆ ต้องมีงานตามเวลาที่ลบแถวที่ `expires_at` ผ่านไปแล้วทิ้ง |
| เปลี่ยน `JWT_SECRET` | เมื่อกุญแจหลุด ต้องเปลี่ยนได้โดยไม่ทำให้ระบบล่ม ซึ่งต้องรองรับกุญแจสองใบพร้อมกันชั่วคราว |
| Log ที่ไม่เผลอบันทึกความลับ | ห้ามบันทึก Token รหัสผ่าน และ `Authorization` ลง Log เด็ดขาด |
| ตรวจความแข็งแรงของรหัสผ่าน | บทนี้ตรวจเพียงความยาว 8 ตัวอักษร ระบบจริงควรเทียบกับรายชื่อรหัสผ่านที่รั่วไหลแล้วด้วย |
| Rate Limit ที่ใช้ที่เก็บกลาง | ตัวนับในบทนี้อยู่ในหน่วยความจำของแต่ละสำเนา เมื่อขยายเป็นหลายสำเนา ต้องย้ายไป Redis |

<aside>
✅

ไม่มีระบบใดปลอดภัยร้อยเปอร์เซ็นต์ สิ่งที่ทำได้คือ **ทำให้ต้นทุนของผู้โจมตีสูงกว่าผลที่เขาจะได้**

และข้อที่สำคัญที่สุดของทั้งวันคือ กฎเรื่องสิทธิ์ต้องอยู่ในที่เดียว ทดสอบได้ และไม่ขึ้นกับว่ามาจาก HTTP หรือไม่

เพราะช่องโหว่ส่วนใหญ่ในงานจริง ไม่ได้เกิดจากการเข้ารหัสที่อ่อนแอ แต่เกิดจาก **การลืมตรวจสิทธิ์ในบางเส้นทาง**

</aside>

## 3.13 Example Problems

โจทย์ทั้งสามข้อทำต่อจากโค้ดที่เสร็จแล้วได้ทันที

### 3.13.1 Endpoint ที่ตรวจสิทธิ์ผิดชั้น

จงพิจารณาโค้ดต่อไปนี้ แล้วตอบคำถามโดยไม่ต้องรัน

```go
func (h Handler) CancelOrder(c *fiber.Ctx) error {
	actor := actorOf(c)
	order, _ := h.uc.Get(actor, c.Params("id"))

	// ตรวจสิทธิ์ไว้ใน Handler
	if order.UserID != actor.ID && actor.Role != "staff" {
		return c.Status(403).JSON(msg("forbidden"))
	}

	return toHTTP(c, h.uc.Cancel(actor, c.Params("id")))
}
```

1. โค้ดนี้ให้ผลลัพธ์ถูกต้องหรือไม่ และผิดหลักข้อใดของ Clean Architecture
2. หากภายหลังเพิ่ม Role ใหม่ชื่อ `manager` ที่ยกเลิกได้ทุกใบ ต้องแก้กี่จุด และมีความเสี่ยงที่จะหลงลืมจุดใด
3. จงเปรียบเทียบกับ `CancelOrder` ที่เขียนไว้ที่ **ขั้นที่ 5** แล้วอธิบายว่าการตรวจสิทธิ์ย้ายไปอยู่ที่ใด

<aside>
⚠️

**Challenge: จงเพิ่ม Role ชื่อ `manager` ที่ยกเลิกได้ทุกใบแต่เปลี่ยนสถานะไม่ได้ เข้าไปในระบบจริง แล้วนับว่าแก้ไปกี่ไฟล์**

</aside>

### 3.13.2 ออเดอร์ของฉัน

จงเพิ่ม Endpoint `GET /api/orders` ที่คืนรายการออเดอร์ของผู้เรียก ตามเงื่อนไขต่อไปนี้

1. ลูกค้าเห็นเฉพาะออเดอร์ของตนเอง ส่วนพนักงานเห็นทุกใบ
2. เพิ่ม Method `ListMine(actor domain.User)` ให้ `OrderUseCase`
3. ขยาย `OrderRepository` แล้วทำให้ทั้ง `memory` และ `postgres` ครบตาม Interface ใหม่
4. ห้ามมีการตรวจ Role อยู่ในชั้น `delivery`

<aside>
⚠️

**Challenge: การกรองควรทำที่ชั้นใด ระหว่างกรองใน SQL ด้วย `WHERE user_id = $1` กับดึงทุกใบขึ้นมาแล้วกรองใน Go จงอธิบายข้อดีข้อเสียของทั้งสองทาง**

</aside>

### 3.13.3 Rate Limit ที่แยกตามผู้ใช้

ตัวนับที่เขียนไว้ในขั้นที่ 11 นับแยกตามหมายเลข IP ซึ่งมีปัญหาเมื่อผู้ใช้หลายคนอยู่หลัง IP เดียวกัน เช่นในเครือข่ายของมหาวิทยาลัย

1. จงแก้ให้เส้นทางที่ผ่าน `RequireAuth` แล้ว นับแยกตาม `actor.ID` แทน IP
2. ส่วนเส้นทางที่ยังไม่รู้ตัวตน เช่นหน้าเข้าสู่ระบบ ให้นับตาม IP ต่อไป
3. อธิบายว่าเหตุใดจึงยังต้องจำกัดตาม IP ที่หน้าเข้าสู่ระบบอยู่

<aside>
🚨

ระวังกับดักหนึ่งข้อ หากนับตาม `actor.ID` เพียงอย่างเดียวที่หน้าเข้าสู่ระบบ จะไม่มีผลใด ๆ เลย

เพราะตอนที่ยังเข้าสู่ระบบไม่สำเร็จ ระบบยังไม่รู้ว่าผู้เรียกเป็นใคร `actor.ID` จึงเป็นค่าว่างเหมือนกันทุกคน

</aside>

## 3.14 Solutions

<aside>
🚨

ให้ลองทำด้วยตนเองจนสุดทางก่อน แล้วจึงเปิดดูเฉลย การอ่านเฉลยที่ถูกต้องไม่ได้ทำให้เขียนโค้ดเองได้

</aside>

### 3.14.1 Endpoint ที่ตรวจสิทธิ์ผิดชั้น

- กดเพื่อดูเฉลย
    
    **ข้อ 1** ให้ผลลัพธ์ถูกต้องในกรณีที่ทดสอบ แต่ผิดหลักสองข้อ
    
    ข้อแรก กฎเรื่องสิทธิ์เป็นกฎทางธุรกิจ จึงต้องอยู่ในชั้น `domain` การเขียนไว้ใน Handler ทำให้ทดสอบไม่ได้หากไม่ส่ง HTTP Request เข้ามา
    
    ข้อที่สอง โค้ดนี้เทียบ `actor.Role != "staff"` ด้วยข้อความดิบ แทนที่จะเรียก `actor.IsStaff()` หากพิมพ์ผิดเป็น `"Staff"` คอมไพเลอร์จะไม่เตือน และทุกคนจะกลายเป็นผู้ไม่มีสิทธิ์เงียบ ๆ
    
    นอกจากนี้ยังละเลย Error จาก `h.uc.Get` หากออเดอร์ไม่มีอยู่จริง `order` จะเป็นค่าว่าง ซึ่ง `order.UserID` เป็นสตริงว่างและไม่ตรงกับผู้ใด จึงคืน `403` ทั้งที่ควรเป็น `404`
    
    **ข้อ 2** ต้องแก้ทุกจุดที่เทียบ Role แบบนี้ ซึ่งกระจายอยู่ตาม Handler ทุกตัว ในระบบจริงมักลืมที่ Endpoint ที่ไม่ค่อยได้ใช้ เช่นหน้าจัดการของผู้ดูแลระบบ
    
    **ข้อ 3** ที่ขั้นที่ 5 Handler ไม่ตัดสินใจเรื่องสิทธิ์เลย เพียงส่ง `actorOf(c)` เข้าไป แล้วให้ `usecase` ถามกฎจาก `domain` การเพิ่ม Role ใหม่จึงแก้ที่ `domain` ที่เดียว
    
    **Challenge** — เพิ่ม Role ชื่อ `manager`
    
    แก้เพียงไฟล์เดียว คือ `domain/user.go` และ `domain/order.go`
    
    ```go
    // domain/user.go
    const (
    	Customer Role = "customer"
    	Staff    Role = "staff"
    	Manager  Role = "manager"
    )
    
    func (u User) IsStaff() bool { return u.Role == Staff }
    func (u User) CanCancelAny() bool { return u.Role == Staff || u.Role == Manager }
    ```
    
    ```go
    // domain/order.go
    func (o Order) CanBeCancelledBy(u User) bool {
    	if !o.CanCancel() {
    		return false
    	}
    	return u.CanCancelAny() || o.OwnedBy(u.ID)
    }
    ```
    
    `CanChangeStatusBy` ไม่ต้องแก้เลย เพราะยังคงเรียก `IsStaff()` ซึ่ง `manager` ไม่เข้าเงื่อนไข ส่วน `delivery` และ `repository` ไม่ถูกแตะแม้แต่บรรทัดเดียว
    

### 3.14.2 ออเดอร์ของฉัน

- กดเพื่อดูเฉลย
    
    ขยาย Interface แล้วเพิ่ม Method ในชั้น Use Case
    
    ```go
    // usecase/order.go
    type OrderRepository interface {
    	ListMenu() ([]domain.MenuItem, error)
    	FindMenuItem(id string) (domain.MenuItem, error)
    	FindOrder(id string) (domain.Order, error)
    	ListOrdersByUser(userID string) ([]domain.Order, error)
    	ListAllOrders() ([]domain.Order, error)
    	Save(o domain.Order) error
    }
    
    // การตัดสินว่าใครเห็นอะไร อยู่ที่ชั้นนี้ ไม่ใช่ที่ Handler
    func (u OrderUseCase) ListMine(actor domain.User) ([]domain.Order, error) {
    	if actor.IsStaff() {
    		return u.repo.ListAllOrders()
    	}
    	return u.repo.ListOrdersByUser(actor.ID)
    }
    ```
    
    Handler ทำเพียงแปลงเข้าและแปลงออก
    
    ```go
    // delivery/http/order.go
    func (h Handler) ListOrders(c *fiber.Ctx) error {
    	orders, err := h.uc.ListMine(actorOf(c))
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
    
    แล้วเพิ่ม Route หนึ่งบรรทัด ในกลุ่มที่ต้องเข้าสู่ระบบก่อน
    
    ```go
    me.Get("/orders", h.ListOrders)
    ```
    
    <aside>
    🚨
    
    การขยาย Interface ครั้งนี้ทำให้มี **4 ที่** ที่ต้องเพิ่ม Method ใหม่ ไม่ใช่ 2 ที่
    
    คือ `memory.OrderRepo`, `postgres.OrderRepo` และอีกสองตัวใน `usecase/fakes_test.go` คือ `stubRepo` กับ `failingRepo`
    
    มิฉะนั้น `go test` จะไม่ผ่าน พร้อมข้อความ `does not implement usecase.OrderRepository`
    
    เป็นราคาเดียวกับที่จ่ายไปแล้วใน Day 2 ขั้นที่ 15
    
    </aside>
    
    **Challenge** — กรองที่ใด
    
    ควรกรองด้วย SQL คือ `WHERE user_id = $1` ด้วยเหตุผลสองข้อ
    
    ข้อแรกคือเรื่องปริมาณ การดึงออเดอร์ทั้งระบบขึ้นมาแล้วค่อยกรองใน Go หมายถึงต้องขนข้อมูลของผู้ใช้ทุกคนผ่านสายมาที่โปรแกรมก่อน เมื่อระบบโตขึ้นจะช้าลงเรื่อย ๆ และเปลืองหน่วยความจำ
    
    ข้อที่สองคือเรื่องความปลอดภัย ข้อมูลของผู้อื่นที่ถูกดึงขึ้นมาไว้ในหน่วยความจำ มีโอกาสหลุดออกไปได้จากความผิดพลาดเล็กน้อย เช่นลืมกรองในเส้นทางใดเส้นทางหนึ่ง หรือเผลอ Log ทั้งก้อนออกไป
    
    หลักที่ใช้คือ **อย่าดึงสิ่งที่ผู้เรียกไม่มีสิทธิ์เห็น ขึ้นมาตั้งแต่แรก**
    

### 3.14.3 Rate Limit ที่แยกตามผู้ใช้

- กดเพื่อดูเฉลย
    
    เปลี่ยนให้ตัวนับรับฟังก์ชันที่บอกว่าจะนับแยกด้วยค่าใด แทนการผูกกับ IP ตายตัว
    
    ```go
    // delivery/http/middleware_ratelimit.go
    type RateLimiter struct {
    	mu     *sync.Mutex
    	hits   map[string]*window
    	limit  int
    	period time.Duration
    	keyOf  func(*fiber.Ctx) string
    }
    
    func NewRateLimiter(limit int, period time.Duration, keyOf func(*fiber.Ctx) string) *RateLimiter {
    	if keyOf == nil {
    		keyOf = func(c *fiber.Ctx) string { return c.IP() }
    	}
    	return &RateLimiter{
    		mu: &sync.Mutex{}, hits: map[string]*window{},
    		limit: limit, period: period, keyOf: keyOf,
    	}
    }
    
    // ในเมธอด Limit เปลี่ยนบรรทัดเดียว
    //   key := c.IP()
    // เป็น
    //   key := rl.keyOf(c)
    ```
    
    แล้วเลือกวิธีนับตอนประกอบระบบใน `main.go`
    
    ```go
    // นับตาม IP สำหรับเส้นทางที่ยังไม่รู้ตัวตน
    loginer := http.NewRateLimiter(5, time.Minute, nil)
    
    // นับตามผู้ใช้ สำหรับเส้นทางที่ผ่าน RequireAuth มาแล้ว
    perUser := http.NewRateLimiter(60, time.Minute, func(c *fiber.Ctx) string {
    	return actorOf(c).ID
    })
    
    me := api.Group("", h.RequireAuth, perUser.Limit)
    ```
    
    <aside>
    🚨
    
    ลำดับสำคัญมาก `perUser.Limit` ต้องอยู่ **หลัง** `h.RequireAuth` เสมอ
    
    หากสลับลำดับ ตัวนับจะทำงานก่อนที่ระบบจะรู้ว่าผู้เรียกเป็นใคร ทำให้ `actorOf(c).ID` เป็นค่าว่างเหมือนกันทุกคน
    
    ผลคือผู้ใช้ทั้งระบบจะใช้ตัวนับก้อนเดียวกัน และเมื่อมีคนใดคนหนึ่งยิงจนเต็มโควตา ทุกคนจะถูกปฏิเสธพร้อมกัน
    
    </aside>
    
    **ข้อ 3** ที่หน้าเข้าสู่ระบบยังต้องจำกัดตาม IP ต่อไป เพราะตอนนั้นระบบยังพิสูจน์ตัวตนไม่สำเร็จ จึงยังไม่มี `actor.ID` ให้ใช้
    
    และเป้าหมายของการจำกัดที่หน้านี้คือกันการไล่เดารหัสผ่าน ซึ่งผู้โจมตีก็ยังไม่มีตัวตนในระบบอยู่แล้ว สิ่งเดียวที่ระบุตัวผู้เรียกได้ในขณะนั้นจึงมีเพียงหมายเลข IP
    

## 3.15 สรุปสิ่งที่ได้ทำวันนี้

| ขั้นที่ | ได้สิ่งใด |
| --- | --- |
| 1 – 2 | `User`, `Role`, Error ของเรื่องสิทธิ์ และกฎทั้งหมดในชั้น `domain` พร้อม Test ที่ไม่ต้องใช้ Token |
| 3 – 4 | `usecase` ประกาศ Interface สามตัว เขียน Register และ Login แล้วส่งตัวตนเข้าไปให้ทุก Method |
| 5 | ชั้น Delivery ทั้งชั้น คือ Middleware, `toHTTP` ที่โตขึ้น และ Handler ของการสมัคร |
| 6 | ค่าลับของระบบยืนยันตัวตน เก็บไว้ใน `.env` |
| 7 – 8 | Adapter ที่ทำงานจริง คือ bcrypt และ JWT พร้อมที่เก็บผู้ใช้ทั้งสองแบบ |
| 9 | ประกอบเข้าด้วยกัน แล้วสมัครสมาชิกครั้งแรก โปรเจกต์กลับมา Build ผ่านทั้งหมด |
| 10 | Migration ที่ 4 และ 5 แล้วย้ายผู้ใช้กับเจ้าของออเดอร์ขึ้น Postgres |
| 11 – 12 | Rate Limiting ที่เขียนขึ้นเอง, Security Header, CORS, Timeout และขนาด Body |
| 13 | ยิงจริงเพื่อพิสูจน์ว่ากฎทำงาน และ Token ปลอมกับ Token หมดอายุใช้ไม่ได้ |
| 14 – 15 | `SessionIssuer` ที่ทำตาม `TokenIssuer` ตัวเดิม พร้อมที่เก็บ Session และ Migration ที่ 6 |
| 16 | ออกจากระบบผ่าน `TokenRevoker` ซึ่งเป็นความสามารถที่ JWT ไม่มี |
| 17 | สลับทั้งระบบระหว่าง JWT กับ Session ด้วยการแก้ `.env` บรรทัดเดียว |

<aside>
💡

ข้อสรุปของทั้งวันคือสิ่งที่พิสูจน์ไปแล้วในขั้นสุดท้าย การเปลี่ยนวิธียืนยันตัวตนทั้งระบบ ใช้การแก้ `.env` เพียงบรรทัดเดียว

เพราะ `domain` รู้จักแค่ `UserID` กับ `Role` ส่วน `usecase` รู้จักแค่ Interface ชื่อ `TokenIssuer` ทั้งสองชั้นไม่เคยเห็นคำว่า JWT หรือ Session เลย

เป็นผลลัพธ์แบบเดียวกับการสลับ Memory ไปเป็น Postgres ใน Day 2 ขั้นที่ 36 คือกฎข้อเดิม เพียงนำมาใช้กับเรื่องใหม่

</aside>

<aside>
⚠️

สิ่งที่ **ยังไม่ได้ทำ** ในวันนี้ และควรทราบว่ายังขาดอยู่

Refresh Token ซึ่งเป็นวิธีที่ระบบจริงใช้รวมข้อดีของทั้งสองแบบเข้าด้วยกัน ดูหัวข้อ 3.11.2

การกวาดล้างแถวใน `sessions` ที่หมดอายุแล้ว ซึ่งยังไม่มีงานตามเวลามาทำให้

การจำกัดจำนวนครั้งที่ยังนับอยู่ในหน่วยความจำ ซึ่งใช้กับหลายสำเนาไม่ได้

HTTPS ซึ่งเป็นสิ่งที่ขาดไม่ได้ในระบบจริง เพราะ Token ที่ส่งผ่าน HTTP ธรรมดา ถูกดักอ่านได้ทั้งเส้นทาง

การบันทึกร่องรอยการเข้าถึง สำหรับตรวจสอบย้อนหลังเมื่อเกิดเหตุ

</aside>