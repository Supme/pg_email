# pg_email
Send email from postgresql

Add to postgresql.conf parameters:
```
email.serverhost = 'smtp.server.tld'
email.serverport = '465'
email.fromname = 'Postgres email sender'
email.fromemail = 'postgres@email.tld'
email.username = 'smtp_user'
email.password ='smtp_password'
```

Get [compiled binary](https://github.com/Supme/pg_email/releases)

Copy files to postgres:
```
cp ./pg_email/pg_email.so /usr/lib/postgresql/11/lib/
cp ./pg_email/pg_email.control /usr/share/postgresql/11/extension/
cp ./pg_email/pg_email--0.1.sql /usr/share/postgresql/11/extension/
```

```
CREATE EXTENSION pg_email;
```

Use in SQL:
```
SELECT SendEmail (to_email, to_name, subject, text_html, text_plain)
```
 return error or blank string then ok
 
 Example:
 ```
 SELECT sendemail('alice@email.tld', 'Alice', 'Email subject from postgresql', '<h1>This HTML</h1><p>Hello!!!</p>', 'This plain text. Hello!');
 ```
