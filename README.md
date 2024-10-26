Assalomu alaykum men twitter clonini tayyorlashga harakat qildim.

1 User Autentifikatsiya:
User ro'yxatdan o'tishi mumkin buning uchun u /api/v1/register endpointiga elektron pochtasi, usernamesi, ismi va parol kiritishi kerak bo'ladi.
User elektron pochtasiga kelgan kodni /api/v1/verify-register endpointiga malumotlarni kirgazishi kerak shunga uning akkaunti yaratiladi
User Login qilishi mumkin buning uchun u /api/v1/login endpointiga usernamesi va parolini kirgazishi kerak bo'ladi shunda u jwt token oladi.

2. Tweetlarni yuklash:
  User Autentifikatsiyadan o'tgandan keyin POST /api/v1/tweet endpinti orqali tweet yuklashi mumkin. Contentga o'z kontent matnini mediaga esa rasm yoki video.

3. Following/Unfollowing User:
  Userlar POST /api/v1/user/{id}/follow endpointi orqali boshqa userga follow qilishlari yoki DELETE /api/v1/user/{id}/unfollow orqali UnFollow qilishlari mumkin.
GET /api/v1/user/{id}/followers va GET /api/v1/user/{id}/followings Endpointlariham bor kim unga follow qilgani va kimlarga follow qilganini bilish uchun.

4. Likes and Views:
  Autentifikatsiyadan o'tgan foydalanuvchilar POST /api/v1/tweet/{id}/like orqali tanlagan tweetiga like bosishi  yoki DELETE /api/v1/tweet/{id}/unlike endpointi orqali unlike qilishlari mumkin .
   PATCH /api/v1/tweet/{id}/views endpointi orqali tweet kurilganlar soni oshiriladi

5. Search:
   Foydalanuvchilar userlarni getlistida searchga sername berish orqali qidirishlari yoki tweetlarni getlistida searchga contentdan bron so'z yoki gapni berish orqali search qilishlari mumkin.

6. Foydalanuvchining tweetiga like bosilganida tweetning egasiga shu odam sizning tweetingizga like bosdi deb habar boradi.

7. Like notificationda kafka message brokeridan foydalanilgan.

8. go test ./... -v  orqali test qilishingiz mumkin

9. https://miro.com/app/board/uXjVLOPN-Nk=/ manzilda system design ni ko'rishingiz mumkin.
10. https://dbdiagram.io/d/Twitter-6718cb4697a66db9a3fb9d15  DB designim.

11. Ro'yxatdan o'tayotganda yuborilgan codni saqlab turishda caching dan foydalanganman
   
