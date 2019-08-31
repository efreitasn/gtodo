package todo

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/data/todo"
	"github.com/efreitasn/gtodo/internal/data/user"
	"github.com/efreitasn/gtodo/internal/utils"
	"github.com/efreitasn/gtodo/pkg/flash"
)

var fetchDoneNotDoneErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while fetching the list of todos.",
}

func (t *Todo) fetchDoneNotDone(w http.ResponseWriter, r *http.Request, tData *template.Data) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userPayload := user.PayloadFromContext(r.Context())

	// Done
	todosDone, err := todo.FetchDone(ctx, t.c, userPayload.ID)

	if err != nil {
		if tData.FlashMessage == nil {
			tData.FlashMessage = fetchDoneNotDoneErrorMsg
		}

		utils.WriteTemplates(w, tData, tData.Mode, "no-todos")

		return
	}

	// Not done
	todosNotDone, err := todo.FetchNotDone(ctx, t.c, userPayload.ID)

	if err != nil {
		if tData.FlashMessage == nil {
			tData.FlashMessage = fetchDoneNotDoneErrorMsg
		}

		utils.WriteTemplates(w, tData, tData.Mode, "no-todos")

		return
	}

	tData.TodosDone = todosDone
	tData.TodosNotDone = todosNotDone

	utils.WriteTemplates(w, tData, tData.Mode, "no-todos")
}
