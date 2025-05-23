package simulation

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"turtle/ctrlApp"
	"turtle/lg"
	"turtle/modelsApp"
)

type SimWorld struct {
	Uid  primitive.ObjectID
	Name string

	SimBehaviours  map[primitive.ObjectID]ISimBehaviour
	SimActors      map[primitive.ObjectID]*SimActor
	SimConnections map[primitive.ObjectID][]ISimBehaviour

	ActorsDefinitions map[primitive.ObjectID]*modelsApp.Actor

	Stepper    SimStepper
	IsOnline   bool
	IdsCounter int64
}

func NewSimWorld() *SimWorld {
	tmp := &SimWorld{}
	tmp.Stepper.End = 100
	tmp.IsOnline = true

	tmp.SimBehaviours = make(map[primitive.ObjectID]ISimBehaviour)
	tmp.SimActors = make(map[primitive.ObjectID]*SimActor)
	tmp.SimConnections = make(map[primitive.ObjectID][]ISimBehaviour)
	tmp.ActorsDefinitions = make(map[primitive.ObjectID]*modelsApp.Actor)

	return tmp
}

func (self *SimWorld) LoadEntities(entities []*modelsApp.Entity) {
	for _, entity := range entities {
		simEntity := SimEntity{}
		simEntity.FromEntity(entity)

		var behaviour ISimBehaviour = NewUndefinedBehaviour()

		entityType := entity.Type

		if entityType == "spawn" {
			behaviour = NewSpawnBehaviour()
		} else if entityType == "process" {
			behaviour = NewProcessBehaviour()
		} else {
			lg.LogE("Unknown entity type [%s]", entityType)
		}

		behaviour.SetWorld(self)
		behaviour.SetEntity(&simEntity)

		self.SimBehaviours[entity.Uid] = behaviour

	}
}

func (self *SimWorld) LoadConnections(connections []*modelsApp.EntityConnection) {
	for _, connection := range connections {

		array, exists := self.SimConnections[connection.A]

		if exists == false {
			array = []ISimBehaviour{}
			self.SimConnections[connection.A] = array
		}

		simBehaviour, found := self.SimBehaviours[connection.B]

		if found {
			self.SimConnections[connection.A] = append(array, simBehaviour)
		} else {
			lg.LogE("Unable to find entity [%s] in world", connection.B.Hex())
		}
	}
}

func (self *SimWorld) PrepareSimulation() {

	for _, behaviour := range self.SimBehaviours {
		behaviour.Init1()
	}

	for _, behaviour := range self.SimBehaviours {
		behaviour.Init2()
	}

}

func (self *SimWorld) GetConnectionsOf(entity primitive.ObjectID) []ISimBehaviour {

	entities, exists := self.SimConnections[entity]

	if exists {
		return entities
	} else {
		return []ISimBehaviour{}
	}

}

func (self *SimWorld) ClearStates() {

}

func (self *SimWorld) Step() {

	lg.LogI(fmt.Sprintf("Step (%d/%d)", self.Stepper.Now, self.Stepper.End))

}

func (self *SimWorld) SpawnActorWithUid(uid primitive.ObjectID) *SimActor {

	definition, exists := self.ActorsDefinitions[uid]

	if exists == false {
		tmp := ctrlApp.GetActor(uid)

		if tmp != nil {
			definition = tmp
			self.ActorsDefinitions[uid] = definition
		} else {
			lg.LogE("SimActor definition [%s] not found", uid.Hex())
			return nil
		}
	}

	actor := SimActor{}
	actor.Id = self.IdsCounter
	self.IdsCounter += 1
	actor.FromActorDefinition(definition)

	return &actor

}
