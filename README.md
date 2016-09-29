# Automate Http/Https Load testing
### status
ไม่สมบูรณ์

โปรแกรมทดสอบ Load ของเว็บแอปพลิเคชันต่างๆ หรือ เอพีไอ ที่ทำงานผ่านโปรโตคอล HTTP, HTTPS
สามารถบอกขีดความสามารถของโปรแกรมที่ท่านพัฒนาในการรองรับจำนวนผู้ใช้งานและข้อมูลต่างๆ
อธิบายข้อมูลจากผลการทดสอบด้วยกราฟชนิดต่างๆ เพื่อให้เห็นความสามารถของแอปพลิเคชันในแต่ละสภาพแวดล้อม
โปรแกรมทำงานอยู่บน wg/wrk ซึ่งเป็น benchmark ที่มีความนิยมในระดับหนึ่ง และเป็นโอเพนซอร์ส ให้ผลลัพธ์, ค่าตัวแปร เช่น
* Request/Second
* Latency
* Data-Transfer/Second
* Socket Error
* Non-2xx Response

### require
* [wg/wrk](https://github.com/wg/wrk)
* [golang](https://golang.org/)

### Available
* Linux
* OSX

### Todo
* Github Webhook
* Bitbucket Webhook