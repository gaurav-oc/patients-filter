import pandas as pd
import os
import psycopg2
import argparse
from configparser import ConfigParser
from sqlalchemy import create_engine
import sys

if __name__ == "__main__":
    iniFormat = """
    [postgresql] # This Line Should be As it is with square Brackets
    host=localhost # Your Database Host Address
    database=dbName # Your Database Name
    user=dbUser # Your Database User
    password=# Your Database Password
    """""
    parser = argparse.ArgumentParser()
    parser.add_argument('-ndif', "--noDatabaseInitializerFile", action="store_true",
                        help="Pass this option if You don't Want to specify The Database Arguments Directly")
    parser.add_argument('-dif', "--databaseInitializerFile", default='database.ini', type=str,
                        help="Database Initializer File")
    parser.add_argument('-f', "--dataFile", default="patients.csv", type=str,
                        help="CSV File Path or Name from data needs to be parsed, default:patients.csv")
    parser.add_argument('-u', '--dbUser', type=str, default="postgres", help="Database User, Default: postgres")
    parser.add_argument('-p', '--dbPassword', type=str, help="Database Password")
    parser.add_argument('-db', '--database', type=str, help="Database Name")
    parser.add_argument('-hs', '--host', type=str, default="localhost",
                        help="Database Host IP or Host Address, Default: localhost")
    parser.add_argument('-t', '--table', type=str, default="patient",
                        help="Table where to store the parsed Results, Default: patient")
    parser.add_argument('--port',type=int,default=5432,help="Port on which Db is running")
    args = parser.parse_args()
    db = {}
    if args.noDatabaseInitializerFile:
        if args.database is not None and args.dbPassword is not None \
                and args.database != "" and args.dbPassword != "":
            db["user"] = args.dbUser
            db["password"] = args.dbPassword
            db["host"] = args.host
            db["database"] = args.database
            db["port"] = args.port
        else:
            print("Please Specify Database Name and Password")
            sys.exit()

    else:
        databaseIniFile = os.path.abspath(args.databaseInitializerFile)

        if os.path.isfile(databaseIniFile):
            parser = ConfigParser()
            section = 'postgresql'
            parser.read(databaseIniFile)
            if parser.has_section(section):
                for item in parser.items(section):
                    db[item[0]] = item[1]
                db["port"] = args.port
            else:

                print("Please Write Your {file} in the Following Format".format(file=args.databaseInitializerFile))
                print(iniFormat)
                sys.exit()
        else:
            print(
                "No File Found. Please Write Your {file} in the Following Format".format(
                    file=args.databaseInitializerFile))
            print(iniFormat)
            sys.exit()
    try:
        print(db["port"])
        datafilePath = os.path.abspath(args.dataFile)
        if os.path.isfile(datafilePath):
            df = pd.read_csv(datafilePath)
            df = df[df.age.notna()]
            df.disease = df.disease.apply(lambda x: x.lower())
            df = df[((df.age > 40) & (df.disease.str.contains('cancer')))]

            conn = psycopg2.connect(**db)
            postgres_url = "postgresql://{user}:{password}@{host}:{port}/{database}".format(**db)
            engine = create_engine(postgres_url)
            df.to_sql('patients', engine, if_exists='append',index=False)
        else:
            print("Record File Doesn't Exist")
    except Exception as e:
        print(db["host"],db["port"])
        print("Exception: ",str(e))
