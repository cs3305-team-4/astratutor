package services

import (
	"errors"

	"github.com/spf13/viper"
	stripe "github.com/stripe/stripe-go/v72"

	stripeAccount "github.com/stripe/stripe-go/v72/account"
	stripeAccountLink "github.com/stripe/stripe-go/v72/accountlink"
	stripeCustomer "github.com/stripe/stripe-go/v72/customer"
	stripeLoginLink "github.com/stripe/stripe-go/v72/loginlink"
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
				MCC:                stripe.String("educational_services"),
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
