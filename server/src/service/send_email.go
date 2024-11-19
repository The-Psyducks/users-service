package service

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(email string, verificationCode int) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("APP_EMAIL_DIR"))
	m.SetHeader("To", email)
  
	m.SetHeader("Subject", "Twitnsap Email verification")
  
	  m.Embed("./twitsnap.png")
  
	  emailBody := fmt.Sprintf(`
		  <div style="font-family: Arial, sans-serif; text-align: center; max-width: 600px; margin: 20px auto; border: 1px solid #ddd; border-radius: 10px; padding: 20px; background-color: #f9f9f9;">
			  <div style="display: flex; align-items: center; justify-content: center; margin-bottom: 20px;">
				  <img src="cid:twitsnap.png" alt="Twitnsap Logo" style="width: 50px; height: 50px; margin-right: 10px;">
				  <h1 style="margin: 0; color: #333;">Twitnsap</h1>
			  </div>
			  <p style="font-size: 16px; color: #555; margin-bottom: 20px;">
				  Verify your email by introducing the verification code below:
			  </p>
			  <div style="margin-top: 20px; font-size: 32px; font-weight: bold; background-color: #f0f0f0; display: inline-block; padding: 20px 30px; border-radius: 5px; color: #333; letter-spacing: 4px;">
				  %d
			  </div>
		  </div>
	  `, verificationCode)
  
	  m.SetBody("text/html", emailBody)
  
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("APP_EMAIL_DIR"), os.Getenv("APP_EMAIL_PASSWORD"))
  
	if err := d.DialAndSend(m); err != nil {
	  return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}