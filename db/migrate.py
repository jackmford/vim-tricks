import csv
import mysql.connector
import os

# Function to create a database connection
def create_connection(host, user, password, database):
    try:
        conn = mysql.connector.connect(
            host=host,
            user=user,
            password=password,
            database=database,
        )
        print("Connected to MySQL database")
        return conn
    except mysql.connector.Error as e:
        print(e)
    return None

# Function to insert data into the database
def insert_data(conn, data):
    sql = ''' INSERT INTO tricks(title, content, lastused)
              VALUES(%s,%s,UTC_TIMESTAMP()) '''
    cursor = conn.cursor()
    cursor.execute(sql, data)
    conn.commit()
    return cursor.lastrowid

# Main function
def main():
    host = 'localhost'
    user = 'vimtricks'
    password = os.environ['DB_PASS']
    database = 'vimtricks'
    csv_file = "tricks.csv"

    # Create a database connection
    conn = create_connection(host, user, password, database)
    if conn:
        # Create table
        # Read data from CSV and insert into database
        with open(csv_file, 'r') as file:
            reader = csv.reader(file)
            for row in reader:
                print(row)
                insert_data(conn, row)
        conn.close()
    else:
        print("Error! Cannot create the database connection.")

if __name__ == '__main__':
    main()
