package domains

import (
	"time"

	"github.com/google/uuid"
)

type ProposalStatus string

const (
	ProposalStatusPending  ProposalStatus = "pending"
	ProposalStatusApproved ProposalStatus = "approved"
	ProposalStatusRejected ProposalStatus = "rejected"
)

type ProductProposal struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	AccountID uuid.UUID
	Status    ProposalStatus
	CreatedAt time.Time
}
