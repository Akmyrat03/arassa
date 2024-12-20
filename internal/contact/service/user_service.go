package service

import (
	"arassachylyk/internal/contact/model"
	"arassachylyk/internal/contact/repository"
	"context"
	"errors"

	"net/mail"
	"net/smtp"
	"regexp"
)

type ContactService struct {
	Repo *repository.ContactRepository
}

func NewContactService(repo *repository.ContactRepository) *ContactService {
	return &ContactService{Repo: repo}
}

// Gmail bilgileri (Uygulama Şifresi gereklidir)
const (
	smtpServer = "smtp.gmail.com"
	smtpPort   = "587"
	smtpEmail  = "akmobile.tm@gmail.com" // Kendi Gmail adresiniz
	smtpPass   = "whclvwobghfdrmqm"      // Gmail uygulama şifresi
)

func (s *ContactService) SendMessage(ctx context.Context, message model.ContactMessage) error {

	if err := validateEmail(message.Email); err != nil {
		return err
	}

	if err := validatePhoneNumber(message.PhoneNumber); err != nil {
		return err
	}

	// 1. Save the message to the database
	if err := s.Repo.SaveMessage(ctx, message); err != nil {
		return err
	}

	// Gmail SMTP setup
	auth := smtp.PlainAuth("", smtpEmail, smtpPass, smtpServer)

	// Email to admin (yourself)
	adminTo := []string{smtpEmail} // Your Gmail
	adminSubject := "Yeni İletişim Mesajı"
	adminBody := "Ad: " + message.Name + "\nEmail: " + message.Email + "\nMesaj:\n" + message.Message
	adminMessage := "From: " + smtpEmail + "\n" +
		"To: " + smtpEmail + "\n" +
		"Subject: " + adminSubject + "\n\n" +
		adminBody

	// Email to the user (sender)
	userTo := []string{message.Email} // User's email
	userSubject := "Mesajınız Alındı"
	userBody := "Salam " + message.Name + ",\n\n" +
		"Hatyňyzy aldyk. Ine, iberen hatyňyz:\n" +
		message.Message + "\n\n" +
		"Sag boluň, \nArassachylyk topary"
	userMessage := "From: " + smtpEmail + "\n" +
		"To: " + message.Email + "\n" +
		"Subject: " + userSubject + "\n\n" +
		userBody

	// 2. Send the emails
	// Send email to admin
	if err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpEmail, adminTo, []byte(adminMessage)); err != nil {
		return err
	}

	// Send email to the user
	if err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpEmail, userTo, []byte(userMessage)); err != nil {
		return err
	}

	return nil
}

// Email validation function
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// Phone number validation function (for Turkmenistan, excluding +993)
func validatePhoneNumber(phoneNumber string) error {
	// Phone number validation regex: valid prefixes (61, 62, 63, 64, 65, 71) followed by 6 digits
	re := regexp.MustCompile(`^(61|62|63|64|65|71)\d{6}$`)
	if !re.MatchString(phoneNumber) {
		return errors.New("geçersiz telefon numarası formatı")
	}
	return nil
}
