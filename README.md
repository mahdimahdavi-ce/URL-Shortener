# URL-Shortener
This is a URL Shortener project which like similar websites you can convert your long URL to a small and handy one. It's awesome!
At first, I only used PostgrSQL database for storing URLs, then I realized that the latency of Postgre is not acceptable for these type of heavy-read websites, so I decided to use Redis in order to cache the frequently used datas and Now It's pretty fast.
Furthermore, I've used Testify in order to write some Unit Test for this project and make sure that It works accurately and is Bug-free.

The technologies that I've used in this project:
- Golang
- GoFiber Framwork (REST API)
- PostgreSQL
- Gorm
- Redis
- Testify (Unit Testing)
