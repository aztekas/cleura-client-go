
def main():
    with open('TODO.md', 'rb') as f:
        data = f.read()
        print(len(data))

if __name__ == "__main__":
    main()
