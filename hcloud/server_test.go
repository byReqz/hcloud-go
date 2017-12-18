package hcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

func TestServerClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerGetResponse{
			Server: schema.Server{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Fatal("no server")
	}
	if server.ID != 1 {
		t.Errorf("unexpected server ID: %v", server.ID)
	}
}

func TestServerClientGetNotFound(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(schema.ErrorResponse{
			Error: schema.Error{
				Code: ErrorCodeNotFound,
			},
		})
	})

	ctx := context.Background()
	server, _, err := env.Client.Server.Get(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		t.Fatal("expected no server")
	}
}

func TestServersList(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		if page := r.URL.Query().Get("page"); page != "2" {
			t.Errorf("expected page 2; got %q", page)
		}
		if perPage := r.URL.Query().Get("per_page"); perPage != "50" {
			t.Errorf("expected per_page 50; got %q", perPage)
		}
		json.NewEncoder(w).Encode(schema.ServerListResponse{
			Servers: []schema.Server{
				{ID: 1},
				{ID: 2},
			},
		})
	})

	opts := ServerListOpts{}
	opts.Page = 2
	opts.PerPage = 50

	ctx := context.Background()
	servers, _, err := env.Client.Server.List(ctx, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 2 {
		t.Fatal("expected 2 servers")
	}
}

func TestServersAll(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Servers []schema.Server `json:"servers"`
			Meta    schema.Meta     `json:"meta"`
		}{
			Servers: []schema.Server{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			Meta: schema.Meta{
				Pagination: &schema.MetaPagination{
					Page:         1,
					LastPage:     1,
					PerPage:      3,
					TotalEntries: 3,
				},
			},
		})
	})

	ctx := context.Background()
	servers, err := env.Client.Server.All(ctx)
	if err != nil {
		t.Fatalf("Servers.List failed: %s", err)
	}
	if len(servers) != 3 {
		t.Fatalf("expected 3 servers; got %d", len(servers))
	}
	if servers[0].ID != 1 || servers[1].ID != 2 || servers[2].ID != 3 {
		t.Errorf("unexpected servers")
	}
}

func TestServersCreate(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"server": {
				"id": 1
			}
		}`)
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.Create(ctx, ServerCreateOpts{
		Name:       "test",
		ServerType: ServerType{ID: 1},
		Image:      Image{ID: 2},
	})
	if err != nil {
		t.Fatalf("Server.Create failed: %s", err)
	}
	if result.Server == nil {
		t.Fatal("no server")
	}
	if result.Server.ID != 1 {
		t.Errorf("unexpected server ID: %v", result.Server.ID)
	}
}

func TestServersDelete(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	ctx := context.Background()
	_, err := env.Client.Server.Delete(ctx, 1)
	if err != nil {
		t.Fatalf("Server.Delete failed: %s", err)
	}
}

func TestServerClientPoweron(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/poweron", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionPoweronResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Poweron(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientReboot(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reboot", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionRebootResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Reboot(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientReset(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reset", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionResetResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Reset(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientShutdown(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/shutdown", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionShutdownResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Shutdown(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientPoweroff(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/poweroff", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionPoweroffResponse{
			Action: schema.Action{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	action, _, err := env.Client.Server.Poweroff(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if action.ID != 1 {
		t.Errorf("unexpected action ID: %d", action.ID)
	}
}

func TestServerClientResetPassword(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/reset_password", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionResetPasswordResponse{
			Action: schema.Action{
				ID: 1,
			},
			RootPassword: "secret",
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.ResetPassword(ctx, &Server{ID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.RootPassword != "secret" {
		t.Errorf("unexpected root password: %v", result.RootPassword)
	}
}

func TestServerClientCreateImageNoOptions(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/create_image", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.ServerActionCreateImageResponse{
			Action: schema.Action{
				ID: 1,
			},
			Image: schema.Image{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	result, _, err := env.Client.Server.CreateImage(ctx, &Server{ID: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.Image.ID != 1 {
		t.Errorf("unexpected image ID: %d", result.Image.ID)
	}
}

func TestServerClientCreateImageWithOptions(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/servers/1/actions/create_image", func(w http.ResponseWriter, r *http.Request) {
		var reqBody schema.ServerActionCreateImageRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatal(err)
		}
		if reqBody.Type == nil || *reqBody.Type != "backup" {
			t.Errorf("unexpected type: %v", reqBody.Type)
		}
		if reqBody.Description == nil || *reqBody.Description != "my backup" {
			t.Errorf("unexpected description: %v", reqBody.Description)
		}
		json.NewEncoder(w).Encode(schema.ServerActionCreateImageResponse{
			Action: schema.Action{
				ID: 1,
			},
			Image: schema.Image{
				ID: 1,
			},
		})
	})

	ctx := context.Background()
	opts := &ServerCreateImageOpts{
		Type:        ImageTypeBackup,
		Description: String("my backup"),
	}
	result, _, err := env.Client.Server.CreateImage(ctx, &Server{ID: 1}, opts)
	if err != nil {
		t.Fatal(err)
	}
	if result.Action.ID != 1 {
		t.Errorf("unexpected action ID: %d", result.Action.ID)
	}
	if result.Image.ID != 1 {
		t.Errorf("unexpected image ID: %d", result.Image.ID)
	}
}
