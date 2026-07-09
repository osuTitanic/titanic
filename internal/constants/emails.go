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

const EmailPasswordChangedSubject = "Your password was changed"
const EmailPasswordChangedBody = `Hi %s,

You are receiving this notification because your account password was changed.
If that was not you, please REPLY IMMEDIATELY and RESET YOUR PASSWORD, as your account may be in danger.

You can reset your password here: %s/account/security

--
Titanic! | %s
`

const EmailAddressChangedSubject = "Your email address was changed"
const EmailAddressChangedBody = `Hi %s,

You are receiving this notification because you (or someone else) changed the email of your account.
If that was not you, please REPLY IMMEDIATELY, as your account may be in danger.

--
Titanic! | %s
`

const EmailReactivateSubject = "Reactivate your account"
const EmailReactivateBody = `Hi %s,

Your account was deactivated, because you have changed your email address.
In order to play again, you will have to re-activate your account, by clicking the link below.

%s/account/verification?id=%d&token=%s

If that was not you, please ignore this email.

--
Titanic! | %s
`
