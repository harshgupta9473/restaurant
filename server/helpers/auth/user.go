package helpers

import (

	"fmt"
	"time"

	"github.com/harshgupta9473/restaurantmanagement/db"
	models "github.com/harshgupta9473/restaurantmanagement/models/user"
)

func GetUserByEmail(email string) (models.User,error){
	db := db.GetDB()
	query:=`SELECT * FROM users WHERE email=$1`
	var user models.User
	err:=db.QueryRow(query,email).Scan(
		&user.Id,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
	);
	
	if err!=nil{
		return user,err
	}
	return user,nil
}

func GetUserByUserId(user_id int64) (models.User,error){
	db := db.GetDB()
	query:=`SELECT * FROM users WHERE id=$1`
	var user models.User
	err:=db.QueryRow(query,user_id).Scan(
		&user.Id,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
	);
	
	if err!=nil{
		return user,err
	}
	return user,nil
}


func CreateNewUser(user models.UserSignup)(int64,error){
	db := db.GetDB()

	query:=`INSERT INTO users
	(first_name,middle_name,last_name,email,password)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id;
	`
	var id int64
	err:=db.QueryRow(query,user.FirstName,user.MiddleName,user.LastName,user.Email,user.Password).Scan(&id)
	if err!=nil{
		return 0,err
	}
	return id,nil
}


func InsertTokenIntoOTPTable(token string,userId int64)error{
	db := db.GetDB()

	query:=`INSERT INTO otp_table
	(user_id,token,expires_at)
	VALUES($1,$2,$3)
	ON CONFLICT (user_id)  -- If user_id already exists
	DO UPDATE SET
	    token=EXCLUDED.token,
		expires_at=EXCLUDED.expires_at;
	`
	expiresAt:=time.Now().Add(15*time.Minute)
	_,err:=db.Exec(query,userId,token,expiresAt)
	return err
}


func FetchOTPDetailsByUserId(user_id int64)(models.OTP,error){
	db:=db.GetDB()

	query:=`SELECT * FROM otp_table 
	WHERE user_id=$1`
	var details models.OTP
	err:=db.QueryRow(query,user_id).Scan(&details.Id,&details.UserId,&details.Token,&details.ExpireAt,&details.Verified)
	
	return details,err
}


func MarkUserVerified(user_id int64)(error){
	db:=db.GetDB()

	tx,err:=db.Begin()
	if err!=nil{
		return fmt.Errorf("failed to start a transaction: %w",err)
	}
	_,err=tx.Exec(`UPDATE users SET verified=TRUE WHERE id=$1`,user_id)
	if err!=nil{
		tx.Rollback()
		return fmt.Errorf("Error updating the users table: %w",err)
	}

	_,err=tx.Exec(`UPDATE otp_table SET verified=TRUE WHERE user_id=$1`,user_id)
	if err!=nil{
		tx.Rollback()
		return fmt.Errorf("Error updating the otp_table: %w",err)
	}
	err=tx.Commit()
	if err!=nil{
		tx.Rollback()
		return fmt.Errorf("error commiting the transaction: %w",err)
	}
	return nil
}