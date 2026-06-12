package constants

const EmailWelcomeSubject = "Welcome to Titanic!"
const EmailWelcomeBody = `Welcome on board, %s!

You've just taken your first steps to experience osu!'s early days.
Before you can play you need to activate your account, by clicking the link below.

%s/account/verification?id=%d&token=%s

We would also recommend joining our discord server:
https://discord.gg/qupv72e7YH

See you in game!

--
Titanic! | %s
`

const EmailPasswordResetSubject = "Reset your password"
const EmailPasswordResetBody = `Hi %s,

You are receiving this notification because you have (or someone pretending to be you has) requested a new password be sent for your account on.
If you did not request this notification, then please ignore it. If you keep receiving it, please contact an administrator.

To use the new password, you need to activate it by clicking the link provided below.

%s/account/verification?id=%d&token=%s&type=password

--
Titanic! | %s
`
