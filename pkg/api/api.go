package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}

type ProfessorHours struct {
	ProfessorID    int     `json:"professor_id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	CommittedHours float64 `json:"committed_hours"`
}

type ScheduleEntry struct {
	DayOfWeek   int    `json:"day_of_week"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	SubjectCode string `json:"subject_code"`
	SubjectName string `json:"subject_name"`
}

type RoomSchedule struct {
	RoomID       int             `json:"room_id"`
	RoomNumber   string          `json:"room_number"`
	BuildingName string          `json:"building_name"`
	Occupied     []ScheduleEntry `json:"occupied_schedules"`
}

// GetProfessorHours godoc
// @Summary Get professor hours
// @Description Returns professor office hours
// @Tags professors
// @Produce json
// @Success 200 {array} api.ProfessorHours
// @Failure 500 {object} map[string]string
// @Router /professor-hours [get]
func (h *Handler) GetProfessorHours(w http.ResponseWriter, r *http.Request) {
	query, err := os.ReadFile("queries/queries.sql")
	if err != nil {
		http.Error(w, "Não foi possível ler o arquivo de query", http.StatusInternalServerError)
		return
	}

	professorQuery := extractQuery(string(query), 1)

	rows, err := h.DB.Query(context.Background(), professorQuery)
	if err != nil {
		log.Printf("Falha na consulta: %v\n", err)
		http.Error(w, "Falha na consulta ao banco", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []ProfessorHours
	for rows.Next() {
		var p ProfessorHours
		if err := rows.Scan(&p.ProfessorID, &p.FirstName, &p.LastName, &p.CommittedHours); err != nil {
			log.Printf("Falha ao ler o registro: %v\n", err)
			continue
		}
		results = append(results, p)
	}

	respondWithJSON(w, http.StatusOK, results)
}

// GetRoomSchedules godoc
// @Summary Get room schedules
// @Description Returns room schedules
// @Tags rooms
// @Produce json
// @Success 200 {array} api.RoomSchedule
// @Failure 500 {object} map[string]string
// @Router /room-schedules [get]
func (h *Handler) GetRoomSchedules(w http.ResponseWriter, r *http.Request) {
	query, err := os.ReadFile("queries/queries.sql")
	if err != nil {
		http.Error(w, "Não foi possível ler o arquivo de queries", http.StatusInternalServerError)
		return
	}

	roomQuery := extractQuery(string(query), 2)

	rows, err := h.DB.Query(context.Background(), roomQuery)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		http.Error(w, "Falha na consulta ao banco", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	roomSchedulesMap := make(map[int]*RoomSchedule)
	for rows.Next() {
		var roomID int
		var roomNumber, buildingName string
		var schedule ScheduleEntry

		if err := rows.Scan(
			&roomID, &roomNumber, &buildingName,
			&schedule.DayOfWeek, &schedule.StartTime, &schedule.EndTime,
			&schedule.SubjectCode, &schedule.SubjectName,
		); err != nil {
			log.Printf("Falha ao ler o registro: %v\n", err)
			continue
		}

		if _, ok := roomSchedulesMap[roomID]; !ok {
			roomSchedulesMap[roomID] = &RoomSchedule{
				RoomID:       roomID,
				RoomNumber:   roomNumber,
				BuildingName: buildingName,
				Occupied:     []ScheduleEntry{},
			}
		}
		roomSchedulesMap[roomID].Occupied = append(roomSchedulesMap[roomID].Occupied, schedule)
	}

	var results []RoomSchedule
	for _, schedule := range roomSchedulesMap {
		results = append(results, *schedule)
	}

	respondWithJSON(w, http.StatusOK, results)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Erro ao serializar a resposta"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func extractQuery(content string, queryNum int) string {
	queries := strings.Split(content, "-- QUERY BREAK --")
	if len(queries) >= queryNum {
		return strings.TrimSpace(queries[queryNum-1])
	}
	return ""
}
