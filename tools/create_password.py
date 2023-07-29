
import hashlib
import bcrypt

def generate(password: str) -> str:
    return bcrypt.hashpw(
        hashlib.md5(password.encode()) \
               .hexdigest() \
               .encode(),
        bcrypt.gensalt()
    ).decode()

if __name__ == "__main__":
    password = input('Please enter your password:\n')
    result = generate(password)

    print(f'This is your hashed password:\n{result}')
