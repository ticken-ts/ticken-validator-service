package fakes

import (
	"github.com/google/uuid"
	"math/big"
	"ticken-validator-service/config"
	"ticken-validator-service/models"
	"ticken-validator-service/repos"
	"ticken-validator-service/utils"
)

type FakeTicketsPopulator struct {
	devUserInfo   config.DevUser
	reposProvider repos.IProvider
}

func NewFakeTicketsPopulator(reposProvider repos.IProvider, devUserInfo config.DevUser) *FakeTicketsPopulator {
	return &FakeTicketsPopulator{
		devUserInfo:   devUserInfo,
		reposProvider: reposProvider,
	}
}

func (populator *FakeTicketsPopulator) Populate() error {
	devEvent, err := populator.createDevEventIfNotExists()
	if err != nil {
		return err
	}
	devAttendant, err := populator.createDevAttendantIfNotExists()
	if err != nil {
		return err
	}

	ticketID := uuid.MustParse("19dd890b-8a17-465f-ab43-0352770c5b9b")

	devTicket := &models.Ticket{
		EventID:             devEvent.EventID,
		TicketID:            ticketID,
		TokenID:             big.NewInt(1),
		ContractAddr:        devEvent.PubBCAddress,
		AttendantID:         devAttendant.AttendantID,
		AttendantWalletAddr: devAttendant.WalletAddress,
	}

	ticketRepo := populator.reposProvider.GetTicketRepository()

	if ticketRepo.FindTicket(devTicket.EventID, devTicket.TicketID) == nil {
		if err := ticketRepo.AddTicket(devTicket); err != nil {
			return err
		}
	}
	return nil
}

func (populator *FakeTicketsPopulator) createDevEventIfNotExists() (*models.Event, error) {
	eventID := uuid.MustParse("60f14a0b-2270-4dd8-90de-752363a0def8")
	organizationID, err := uuid.Parse(populator.devUserInfo.OrganizationID)
	if err != nil {
		return nil, err
	}

	devEvent := &models.Event{
		EventID:        eventID,
		OrganizerID:    uuid.New(),
		PvtBCChannel:   "pvtbc-fake-channel",
		PubBCAddress:   "pubbc-fake-address",
		OrganizationID: organizationID,
	}

	eventRepo := populator.reposProvider.GetEventRepository()

	if !eventRepo.AnyWithID(eventID) {
		if err := eventRepo.AddEvent(devEvent); err != nil {
			return nil, err
		}
	}

	return devEvent, nil
}

func (populator *FakeTicketsPopulator) createDevAttendantIfNotExists() (*models.Attendant, error) {
	attendantID := uuid.MustParse("19dd890b-8a17-465f-ab43-0352770c5b9b")

	attendantPubKey :=
		"-----BEGIN RSA PUBLIC KEY-----" +
			"\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDUqRecbJUp6pVi5cCpB97cxNav" +
			"\nHxjlPCa4sB/MwcGbo/nuP/iPJf6X4XeRM1xFey5b3STSbhBKmM01+WR6XmUBvzdE" +
			"\nPntm/3Byn9OfKIZcXJoRP1xpFbwfBcSlpZj0bVX0i+1asWWCvdykjcyNy3CMYCSR" +
			"\nHzquQ87BNrh7g7nfjwIDAQAB" +
			"\n-----END RSA PUBLIC KEY-----"

	attendantPrivKey :=
		"-----BEGIN RSA PRIVATE KEY-----" +
			"\nMIICWwIBAAKBgQDUqRecbJUp6pVi5cCpB97cxNavHxjlPCa4sB/MwcGbo/nuP/iP" +
			"\nJf6X4XeRM1xFey5b3STSbhBKmM01+WR6XmUBvzdEPntm/3Byn9OfKIZcXJoRP1xp" +
			"\nFbwfBcSlpZj0bVX0i+1asWWCvdykjcyNy3CMYCSRHzquQ87BNrh7g7nfjwIDAQAB" +
			"\nAoGAW3URM3O7PtilQHAgyFEbNoTs80mDcmrJGFqegne9pQsDXMRkSGQFtxn/SxH0" +
			"\nl+kfCeD0ig9NsFdAwfqsjLf15d6KvXZOlV8zeHHB2qLeW+1orNmThYurAPo4+MdU" +
			"\nEdpYe3bGB8ZidEpDPrg1zhDlePnN9bsKSFOlCzXJRKGBTekCQQD0AiTTF0dNbRmG" +
			"\nzGRmmJoLDVW32RfNCLg4CnGFIuumAKntglKYNig4ke9hNTd5+sSD/jdWSBnTTYRl" +
			"\nhnCZpNI1AkEA3xyQC93M7zb4FPKEXhAaZ/XyJGgarRtNr//qGVBypP85Yucd0h5j" +
			"\nyamh+ArJZouplZX0s7Pi/2bhZl4ne2PjMwJAShDghqa5QPpN1knya+YEVDh+/WhL" +
			"\nPjRYXsJkxOndp6zp56s4UPWXbdx2UgZqSX9h6ULgHzORiz8rYfnV8f1CxQJAWoFC" +
			"\nqZ2i2VMKFa0/Js0PeSaawEv+rkQKIqAEfZpVtzrVM5qfTTIItrB6RJ1Tj6aN92Eq" +
			"\nL4+EQKiiPJ1rFLGzYwJAGpZsgLpoJq6NppNTuuKK1OJj0NSoeIdYl0UftG5pG1V+" +
			"\nHw09Kn47X96bdsycEaKN/aZUPGVw9Bqtv/gsw/FTQQ==" +
			"\n-----END RSA PRIVATE KEY-----"

	_, err := utils.LoadRSA(attendantPrivKey, attendantPubKey)
	if err != nil {
		return nil, err
	}

	devAttendant := &models.Attendant{
		AttendantID:   attendantID,
		WalletAddress: "fake-wallet-addr",
		PublicKey:     []byte(attendantPubKey),
	}

	attendantRepo := populator.reposProvider.GetAttendantRepository()

	if !attendantRepo.AnyWithID(attendantID) {
		if err := attendantRepo.AddAttendant(devAttendant); err != nil {
			return nil, err
		}
	}

	return devAttendant, nil
}
