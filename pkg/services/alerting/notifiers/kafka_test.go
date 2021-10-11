package notifiers

import (
	"testing"

	"github.com/grafana/grafana/pkg/services/secrets/database"
	secretsManager "github.com/grafana/grafana/pkg/services/secrets/manager"
	"github.com/grafana/grafana/pkg/services/sqlstore"

	"github.com/grafana/grafana/pkg/components/simplejson"
	"github.com/grafana/grafana/pkg/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestKafkaNotifier(t *testing.T) {
	Convey("Kafka notifier tests", t, func() {
		store := database.ProvideSecretsStore(sqlstore.InitTestDB(t))
		secretsService := secretsManager.SetupTestService(t, store)
		Convey("Parsing alert notification from settings", func() {
			Convey("empty settings should return error", func() {
				json := `{ }`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &models.AlertNotification{
					Name:     "kafka_testing",
					Type:     "kafka",
					Settings: settingsJSON,
				}

				_, err := NewKafkaNotifier(model, secretsService.GetDecryptedValue)
				So(err, ShouldNotBeNil)
			})

			Convey("settings should send an event to kafka", func() {
				json := `
				{
					"kafkaRestProxy": "http://localhost:8082",
					"kafkaTopic": "topic1"
				}`

				settingsJSON, _ := simplejson.NewJson([]byte(json))
				model := &models.AlertNotification{
					Name:     "kafka_testing",
					Type:     "kafka",
					Settings: settingsJSON,
				}

				not, err := NewKafkaNotifier(model, secretsService.GetDecryptedValue)
				kafkaNotifier := not.(*KafkaNotifier)

				So(err, ShouldBeNil)
				So(kafkaNotifier.Name, ShouldEqual, "kafka_testing")
				So(kafkaNotifier.Type, ShouldEqual, "kafka")
				So(kafkaNotifier.Endpoint, ShouldEqual, "http://localhost:8082")
				So(kafkaNotifier.Topic, ShouldEqual, "topic1")
			})
		})
	})
}
