package market

import (
	"fmt"

	"github.com/lovyou-ai/eventgraph/go/pkg/event"
	"github.com/lovyou-ai/eventgraph/go/pkg/store"
	"github.com/lovyou-ai/eventgraph/go/pkg/types"
)

// Endorsement holds the data from a market.reputation.endorsement event.
type Endorsement struct {
	ID         types.EventID
	EndorserID types.ActorID
	SubjectID  types.ActorID
	Skill      string
	Workspace  string
}

// Reputation represents the accumulated endorsements for an actor, keyed by skill.
type Reputation struct {
	SubjectID   types.ActorID
	SkillCounts map[string]int
}

// Review holds the data from a market.reputation.review event.
type Review struct {
	ID         types.EventID
	ReviewerID types.ActorID
	SubjectID  types.ActorID
	TaskID     types.EventID
	Rating     int
	Note       string
	Workspace  string
}

// ReputationStore creates and queries reputation events as auditable events on the shared graph.
type ReputationStore struct {
	store   store.Store
	factory *event.EventFactory
	signer  event.Signer
}

// NewReputationStore creates a new ReputationStore backed by the given event store.
func NewReputationStore(s store.Store, factory *event.EventFactory, signer event.Signer) *ReputationStore {
	return &ReputationStore{store: s, factory: factory, signer: signer}
}

// AddEndorsement records a market.reputation.endorsement event on the graph.
// The endorser attests that the subject has demonstrated the named skill.
func (rs *ReputationStore) AddEndorsement(
	endorser types.ActorID,
	subject types.ActorID,
	skill string,
	causes []types.EventID,
	convID types.ConversationID,
) (Endorsement, error) {
	if skill == "" {
		return Endorsement{}, fmt.Errorf("skill is required")
	}
	content := EndorsementContent{
		EndorserID: endorser,
		SubjectID:  subject,
		Skill:      skill,
	}
	ev, err := rs.factory.Create(EventTypeEndorsement, endorser, content, causes, convID, rs.store, rs.signer)
	if err != nil {
		return Endorsement{}, fmt.Errorf("create endorsement event: %w", err)
	}
	stored, err := rs.store.Append(ev)
	if err != nil {
		return Endorsement{}, fmt.Errorf("append endorsement event: %w", err)
	}
	return Endorsement{
		ID:         stored.ID(),
		EndorserID: endorser,
		SubjectID:  subject,
		Skill:      skill,
		Workspace:  content.Workspace,
	}, nil
}

// GetReputation returns the accumulated endorsement counts per skill for the given actor.
// Returns an empty SkillCounts map if the actor has no endorsements.
func (rs *ReputationStore) GetReputation(subject types.ActorID) (Reputation, error) {
	page, err := rs.store.ByType(EventTypeEndorsement, 1000, types.None[types.Cursor]())
	if err != nil {
		return Reputation{}, fmt.Errorf("fetch endorsement events: %w", err)
	}
	counts := make(map[string]int)
	for _, ev := range page.Items() {
		c, ok := ev.Content().(EndorsementContent)
		if !ok || c.SubjectID != subject {
			continue
		}
		counts[c.Skill]++
	}
	return Reputation{
		SubjectID:   subject,
		SkillCounts: counts,
	}, nil
}

// AddReview records a market.reputation.review event on the graph.
// rating must be between 1 and 5. taskID may be zero if not tied to a specific task.
func (rs *ReputationStore) AddReview(
	reviewer types.ActorID,
	subject types.ActorID,
	taskID types.EventID,
	rating int,
	note string,
	causes []types.EventID,
	convID types.ConversationID,
) (Review, error) {
	if rating < 1 || rating > 5 {
		return Review{}, fmt.Errorf("rating must be between 1 and 5")
	}
	content := ReviewContent{
		ReviewerID: reviewer,
		SubjectID:  subject,
		TaskID:     taskID,
		Rating:     rating,
		Note:       note,
	}
	ev, err := rs.factory.Create(EventTypeReview, reviewer, content, causes, convID, rs.store, rs.signer)
	if err != nil {
		return Review{}, fmt.Errorf("create review event: %w", err)
	}
	stored, err := rs.store.Append(ev)
	if err != nil {
		return Review{}, fmt.Errorf("append review event: %w", err)
	}
	return Review{
		ID:         stored.ID(),
		ReviewerID: reviewer,
		SubjectID:  subject,
		TaskID:     taskID,
		Rating:     rating,
		Note:       note,
		Workspace:  content.Workspace,
	}, nil
}

// ListReviews returns all reviews for the given subject actor.
func (rs *ReputationStore) ListReviews(subject types.ActorID) ([]Review, error) {
	page, err := rs.store.ByType(EventTypeReview, 1000, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("fetch review events: %w", err)
	}
	var reviews []Review
	for _, ev := range page.Items() {
		c, ok := ev.Content().(ReviewContent)
		if !ok || c.SubjectID != subject {
			continue
		}
		reviews = append(reviews, Review{
			ID:         ev.ID(),
			ReviewerID: c.ReviewerID,
			SubjectID:  c.SubjectID,
			TaskID:     c.TaskID,
			Rating:     c.Rating,
			Note:       c.Note,
			Workspace:  c.Workspace,
		})
	}
	return reviews, nil
}

// AddEndorsementInWorkspace records a market.reputation.endorsement event scoped to a workspace.
func (rs *ReputationStore) AddEndorsementInWorkspace(
	endorser types.ActorID,
	subject types.ActorID,
	skill string,
	workspace string,
	causes []types.EventID,
	convID types.ConversationID,
) (Endorsement, error) {
	if skill == "" {
		return Endorsement{}, fmt.Errorf("skill is required")
	}
	if workspace == "" {
		return Endorsement{}, fmt.Errorf("workspace is required")
	}
	content := EndorsementContent{
		EndorserID: endorser,
		SubjectID:  subject,
		Skill:      skill,
		Workspace:  workspace,
	}
	ev, err := rs.factory.Create(EventTypeEndorsement, endorser, content, causes, convID, rs.store, rs.signer)
	if err != nil {
		return Endorsement{}, fmt.Errorf("create endorsement event: %w", err)
	}
	stored, err := rs.store.Append(ev)
	if err != nil {
		return Endorsement{}, fmt.Errorf("append endorsement event: %w", err)
	}
	return Endorsement{
		ID:         stored.ID(),
		EndorserID: endorser,
		SubjectID:  subject,
		Skill:      skill,
		Workspace:  workspace,
	}, nil
}

// AddReviewInWorkspace records a market.reputation.review event scoped to a workspace.
func (rs *ReputationStore) AddReviewInWorkspace(
	reviewer types.ActorID,
	subject types.ActorID,
	taskID types.EventID,
	rating int,
	note string,
	workspace string,
	causes []types.EventID,
	convID types.ConversationID,
) (Review, error) {
	if rating < 1 || rating > 5 {
		return Review{}, fmt.Errorf("rating must be between 1 and 5")
	}
	if workspace == "" {
		return Review{}, fmt.Errorf("workspace is required")
	}
	content := ReviewContent{
		ReviewerID: reviewer,
		SubjectID:  subject,
		TaskID:     taskID,
		Rating:     rating,
		Note:       note,
		Workspace:  workspace,
	}
	ev, err := rs.factory.Create(EventTypeReview, reviewer, content, causes, convID, rs.store, rs.signer)
	if err != nil {
		return Review{}, fmt.Errorf("create review event: %w", err)
	}
	stored, err := rs.store.Append(ev)
	if err != nil {
		return Review{}, fmt.Errorf("append review event: %w", err)
	}
	return Review{
		ID:         stored.ID(),
		ReviewerID: reviewer,
		SubjectID:  subject,
		TaskID:     taskID,
		Rating:     rating,
		Note:       note,
		Workspace:  workspace,
	}, nil
}

// ListEndorsementsByWorkspace returns all endorsements scoped to the given workspace.
func (rs *ReputationStore) ListEndorsementsByWorkspace(workspace string) ([]Endorsement, error) {
	page, err := rs.store.ByType(EventTypeEndorsement, 1000, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("fetch endorsement events: %w", err)
	}
	var endorsements []Endorsement
	for _, ev := range page.Items() {
		c, ok := ev.Content().(EndorsementContent)
		if !ok || c.Workspace != workspace {
			continue
		}
		endorsements = append(endorsements, Endorsement{
			ID:         ev.ID(),
			EndorserID: c.EndorserID,
			SubjectID:  c.SubjectID,
			Skill:      c.Skill,
			Workspace:  c.Workspace,
		})
	}
	return endorsements, nil
}

// ListReviewsByWorkspace returns all reviews scoped to the given workspace.
func (rs *ReputationStore) ListReviewsByWorkspace(workspace string) ([]Review, error) {
	page, err := rs.store.ByType(EventTypeReview, 1000, types.None[types.Cursor]())
	if err != nil {
		return nil, fmt.Errorf("fetch review events: %w", err)
	}
	var reviews []Review
	for _, ev := range page.Items() {
		c, ok := ev.Content().(ReviewContent)
		if !ok || c.Workspace != workspace {
			continue
		}
		reviews = append(reviews, Review{
			ID:         ev.ID(),
			ReviewerID: c.ReviewerID,
			SubjectID:  c.SubjectID,
			TaskID:     c.TaskID,
			Rating:     c.Rating,
			Note:       c.Note,
			Workspace:  c.Workspace,
		})
	}
	return reviews, nil
}

// GetReputationByWorkspace returns accumulated endorsement counts per skill for the given actor,
// filtered to only endorsements within the named workspace.
func (rs *ReputationStore) GetReputationByWorkspace(subject types.ActorID, workspace string) (Reputation, error) {
	endorsements, err := rs.ListEndorsementsByWorkspace(workspace)
	if err != nil {
		return Reputation{}, fmt.Errorf("fetch workspace endorsements: %w", err)
	}
	counts := make(map[string]int)
	for _, e := range endorsements {
		if e.SubjectID == subject {
			counts[e.Skill]++
		}
	}
	return Reputation{
		SubjectID:   subject,
		SkillCounts: counts,
	}, nil
}
