package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	stripe "github.com/stripe/stripe-go/v72"
	"gorm.io/gorm"

	stripeAccount "github.com/stripe/stripe-go/v72/account"
	stripeAccountLink "github.com/stripe/stripe-go/v72/accountlink"
	stripeCheckoutSession "github.com/stripe/stripe-go/v72/checkout/session"
	stripeCustomer "github.com/stripe/stripe-go/v72/customer"
	stripeLoginLink "github.com/stripe/stripe-go/v72/loginlink"
	stripePaymentIntent "github.com/stripe/stripe-go/v72/paymentintent"
	stripePaymentMethod "github.com/stripe/stripe-go/v72/paymentmethod"
	stripePayout "github.com/stripe/stripe-go/v72/payout"
	stripeRefund "github.com/stripe/stripe-go/v72/refund"
	stripeTransfer "github.com/stripe/stripe-go/v72/transfer"
)

// SetupBilling sets up billing for an account
// The account must have a valid profile
func (ac *Account) SetupBilling() error {
	if ac.Profile == nil {
		return errors.New("cannot setup billing on an account without a profile")
	}

	if ac.Type == Tutor {
		billAcc, err := stripeAccount.New(&stripe.AccountParams{
			Individual: &stripe.PersonParams{
				Email: stripe.String(ac.Email),
			},
			Email: stripe.String(ac.Email),
			Type:  stripe.String("express"),
			BusinessProfile: &stripe.AccountBusinessProfileParams{
				MCC:                stripe.String("8299"),
				ProductDescription: stripe.String("Tutor for AstraTutor"),
				SupportEmail:       stripe.String(ac.Email),
				Name:               stripe.String("AstaTutor - " + ac.Profile.FirstName + " " + ac.Profile.LastName),
			},
			Capabilities: &stripe.AccountCapabilitiesParams{
				Transfers: &stripe.AccountCapabilitiesTransfersParams{
					Requested: stripe.Bool(true),
				},
				CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
					Requested: stripe.Bool(true),
				},
				SEPADebitPayments: &stripe.AccountCapabilitiesSEPADebitPaymentsParams{
					Requested: stripe.Bool(true),
				},
			},
			BusinessType: stripe.String("individual"),
			Settings: &stripe.AccountSettingsParams{
				Payouts: &stripe.AccountSettingsPayoutsParams{
					Schedule: &stripe.PayoutScheduleParams{
						Interval: stripe.String("manual"),
					},
					StatementDescriptor: stripe.String("AstraTutor"),
				},
			},
		})

		if err != nil {
			return err
		}

		ac.StripeID = billAcc.ID
	} else if ac.Type == Student {
		cusmAcc, err := stripeCustomer.New(&stripe.CustomerParams{
			Name:  stripe.String(ac.Profile.FirstName + " " + ac.Profile.LastName),
			Email: stripe.String(ac.Email),
		})
		if err != nil {
			return err
		}

		ac.StripeID = cusmAcc.ID
	}

	return nil
}

func (ac *Account) IsTutorBillingOnboarded() (bool, error) {
	billAcc, err := stripeAccount.GetByID(ac.StripeID, nil)
	if err != nil {
		return false, err
	}

	return billAcc.ChargesEnabled, nil
}

func (ac *Account) IsTutorBillingRequirementsMet() (bool, error) {
	billAcc, err := stripeAccount.GetByID(ac.StripeID, nil)
	if err != nil {
		return false, err
	}

	return len(billAcc.Requirements.CurrentlyDue) == 0 && len(billAcc.Requirements.EventuallyDue) == 0, nil
}

func (ac *Account) GetTutorBillingOnboardURL() (string, error) {
	link, err := stripeAccountLink.New(&stripe.AccountLinkParams{
		Account:    stripe.String(ac.StripeID),
		RefreshURL: stripe.String(viper.GetString("billing.stripe.account_link.refresh_url")),
		ReturnURL:  stripe.String(viper.GetString("billing.stripe.account_link.return_url")),
		Type:       stripe.String("account_onboarding"),
	})
	if err != nil {
		return "", err
	}

	return link.URL, nil
}

// GetTutorBillingPanelURL returns a link that a user can use to access their billing account on Stripe
func (ac *Account) GetTutorBillingPanelURL() (string, error) {
	params := &stripe.LoginLinkParams{
		Account: stripe.String(ac.StripeID),
	}

	ll, err := stripeLoginLink.New(params)
	if err != nil {
		return "", err
	}

	return ll.URL, nil
}

func (l *Lesson) SetupPaymentIntent() error {
	subjectTaught := l.SubjectTaught
	student := l.Student

	intent, err := stripePaymentIntent.New(&stripe.PaymentIntentParams{
		Amount:   stripe.Int64(subjectTaught.Price),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		SetupFutureUsage: stripe.String(string(stripe.PaymentIntentSetupFutureUsageOffSession)),
		Customer:         &student.StripeID,
	})
	if err != nil {
		return err
	}

	l.PaymentIntentID = intent.ID
	l.PayoutAmount = ((subjectTaught.Price) / 100) * (100 - viper.GetInt64("billing.profit_margin"))
	l.PriceAmount = (subjectTaught.Price)
	return nil
}

func (acc *Account) CreateCardSetupSession(successPath string, cancelPath string) (string, error) {
	types := []*string{stripe.String("card")}
	checkout, err := stripeCheckoutSession.New(&stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String(viper.GetString("ui.base_url") + successPath),
		CancelURL:          stripe.String(viper.GetString("ui.base_url") + cancelPath),
		PaymentMethodTypes: types,
		Customer:           &acc.StripeID,
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSetup)),
	})
	if err != nil {
		return "", err
	}

	return checkout.ID, err
}

func (acc *Account) GetCards() ([]stripe.PaymentMethod, error) {
	res := stripePaymentMethod.List(&stripe.PaymentMethodListParams{
		Customer: stripe.String(acc.StripeID),
		Type:     stripe.String("card"),
	})

	ret := []stripe.PaymentMethod{}
	for res.Next() {
		ret = append(ret, *res.PaymentMethod())
	}

	return ret, nil
}

func (acc *Account) DeleteCard(id string) error {
	pm, err := stripePaymentMethod.Get(id, nil)
	if err != nil {
		return err
	}

	if pm.Customer.ID != acc.StripeID {
		return errors.New("this payment method is not associated with the specified account")
	}

	_, err = stripePaymentMethod.Detach(id, nil)
	if err != nil {
		return err
	}

	return nil
}

type PayeePayment struct {
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Amount      int64     `json:"amount"`
	Remarks     string    `json:"remarks"`
}

// A PayerPayment is a payment that an account has recieved that can be payed out to a bank account
type PayerPayment struct {
	// Description of the payment
	Description string `json:"description"`

	// Date
	Date time.Time `json:"date"`

	// Amount
	Amount int64 `json:"amount"`

	// Remarks about the payment
	Remarks string `json:"remarks"`

	// AvailableForPayout true if can be paid out
	AvailableForPayout bool `json:"available_for_payout"`

	// PaidOut is true if the payment has been paid out
	PaidOut bool `json:"paid_out"`
}

// PayoutInfo concerns info about how much money the account can be paid out
type PayoutInfo struct {
	// PayoutBalance is the amount of money that can be paid out (in cents)
	PayoutBalance int64 `json:"payout_balance"`
}

func (acc *Account) GetPayeesPayments() ([]PayeePayment, error) {
	if acc.Type == Tutor {
		return nil, errors.New("tutors cannot make payments, they do not pay accounts")
	}

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	// Find every lesson the student has that has been paid for
	var lessons []Lesson
	err = db.Where(&Lesson{
		StudentID: acc.ID,
		Paid:      true,
	}).Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	var payees []PayeePayment
	payees = []PayeePayment{}

	for _, lesson := range lessons {
		payees = append(payees, PayeePayment{
			Description: lesson.StartTime.Format("Lesson on 2006-01-02"),
			Date:        *lesson.DatePaid,
			Amount:      lesson.PriceAmount,
			Remarks:     "",
		})
	}

	return payees, nil
}

// GetPayoutInfo returns information about payouts
func (acc *Account) GetPayoutInfo() (*PayoutInfo, error) {
	payers, err := acc.GetPayersPayments()
	if err != nil {
		return nil, err
	}

	var cents int64
	cents = 0

	for _, payer := range payers {
		if payer.AvailableForPayout && !payer.PaidOut {
			cents += payer.Amount
		}
	}

	return &PayoutInfo{
		PayoutBalance: cents,
	}, nil
}

func (acc *Account) Payout() error {
	if acc.Type != Tutor {
		return errors.New("only tutors can receive payouts")
	}

	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// Find every lesson the tutor has that has been paid for
		var unpaidOutlessons []Lesson
		err := tx.Where(&Lesson{
			TutorID: acc.ID,
			Paid:    true,
			PaidOut: false,
		}, "TutorID", "Paid", "PaidOut").Find(&unpaidOutlessons).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		var amount int64
		amount = 0
		paidLessonIds := []uuid.UUID{}

		for _, lesson := range unpaidOutlessons {
			if !viper.GetBool("billing.allow_instant_payouts") {
				duration := time.Now().Sub(lesson.StartTime)
				numDays := int(duration.Hours()) / 24

				if numDays > 14 {
					amount += lesson.PayoutAmount
					paidLessonIds = append(paidLessonIds, lesson.ID)
				}
			} else {
				amount += lesson.PayoutAmount
				paidLessonIds = append(paidLessonIds, lesson.ID)
			}
		}

		now := time.Now()

		err = tx.Model(Lesson{}).Where("id IN ?", paidLessonIds).Updates(Lesson{PaidOut: true, DatePaidOut: &now}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		transferParams := &stripe.TransferParams{
			Amount:      stripe.Int64(amount),
			Currency:    stripe.String(string(stripe.CurrencyEUR)),
			Destination: stripe.String(acc.StripeID),
		}

		_, err = stripeTransfer.New(transferParams)
		if err != nil {
			tx.Rollback()
			return err
		}

		params := &stripe.PayoutParams{
			Amount:   stripe.Int64(amount),
			Currency: stripe.String(string(stripe.CurrencyEUR)),
		}
		params.SetStripeAccount(acc.StripeID)

		_, err = stripePayout.New(params)
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

// GetPayersPayments returns a list of payments that the account has received
func (acc *Account) GetPayersPayments() ([]PayerPayment, error) {
	if acc.Type == Student {
		return nil, errors.New("students cannot receive funds, they do not have payers")
	}

	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	// Find every lesson the tutor has that has been paid for
	var lessons []Lesson
	err = db.Where(&Lesson{
		TutorID: acc.ID,
		Paid:    true,
	}).Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	var payers []PayerPayment
	payers = []PayerPayment{}

	for _, lesson := range lessons {
		if lesson.PaidOut {
			payers = append(payers, PayerPayment{
				Description:        lesson.StartTime.Format("Lesson on 2006-01-02"),
				Date:               *lesson.DatePaid,
				Amount:             lesson.PayoutAmount,
				Remarks:            "",
				AvailableForPayout: false,
				PaidOut:            true,
			})

			continue
		}

		if !viper.GetBool("billing.allow_instant_payouts") {
			duration := time.Now().Sub(lesson.StartTime)
			numDays := int(duration.Hours()) / 24

			if numDays >= 14 {
				payers = append(payers, PayerPayment{
					Description:        lesson.StartTime.Format("Lesson on 2006-01-02"),
					Date:               *lesson.DatePaid,
					Amount:             lesson.PayoutAmount,
					Remarks:            "",
					AvailableForPayout: true,
					PaidOut:            false,
				})
			} else {
				payers = append(payers, PayerPayment{
					Description:        lesson.StartTime.Format("Lesson on 2006-01-02"),
					Date:               *lesson.DatePaid,
					Amount:             lesson.PayoutAmount,
					Remarks:            fmt.Sprintf("Available for payout in %d days", (14 - numDays)),
					AvailableForPayout: false,
					PaidOut:            false,
				})
			}
		} else {
			payers = append(payers, PayerPayment{
				Description:        lesson.StartTime.Format("Lesson on 2006-01-02"),
				Date:               *lesson.DatePaid,
				Amount:             lesson.PayoutAmount,
				Remarks:            "",
				AvailableForPayout: true,
				PaidOut:            false,
			})
		}
	}

	return payers, nil
}

func (l *Lesson) Refund() error {
	if l.Refunded == true {
		return nil
	}

	_, err := stripeRefund.New(&stripe.RefundParams{
		PaymentIntent: stripe.String(l.PaymentIntentID),
	})
	if err != nil {
		return err
	}
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&l).Updates(&Lesson{
			Refunded: true,
		}).Error
	})

	if err != nil {
		return err
	}

	return nil
}

// RereshPaidStatus double checks with Stripe if the lesson has been paid for yet, and if it has, updates the lesson
func (l *Lesson) RefreshPaidStatus() error {
	if l.Paid == true {
		return nil
	}

	intent, err := stripePaymentIntent.Get(l.PaymentIntentID, nil)
	if err != nil {
		return err
	}

	if intent.Status == stripe.PaymentIntentStatusSucceeded {
		db, err := database.Open()
		if err != nil {
			return err
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			now := time.Now()

			return tx.Model(l).Updates(&Lesson{
				Paid:     true,
				DatePaid: &now,
			}).Error
		})
		if err != nil {
			return err
		}

		l.Paid = true
	}

	return nil
}

func (l *Lesson) GetPaymentIntentClientSecret() (string, error) {
	intent, err := stripePaymentIntent.Get(l.PaymentIntentID, nil)
	if err != nil {
		return "", err
	}

	return intent.ClientSecret, err
}
