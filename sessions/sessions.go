package sessions

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"
	"task-manager/models"
	st "task-manager/store"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func CreateUserSession(store *session.Store, c *fiber.Ctx, db *sql.DB, username string) error {
	// Get or create session
	s, _ := store.Get(c)
	// fmt.Println(s.Fresh())

	// If this is a new session
	if s.Fresh() {
		// Get session ID
		sid := s.ID()

		//Get user ID
		// uid := c.Params("uid")

		// Save session data
		s.Set("username", username)
		s.Set("sid", sid)
		s.Set("ip", c.Context().RemoteIP().String())
		s.Set("login", time.Unix(time.Now().Unix(), 0).UTC().String())
		s.Set("ua", string(c.Request().Header.UserAgent()))
		s.SetExpiry(time.Second * 30)

		log.Print(username)
		err := s.Save()
		if err != nil {
			// log.Println(err)
			return err
		}

		// Save user reference
		stmt, err := db.Prepare(`UPDATE session_store SET u = $1 WHERE k = $2`)
		if err != nil {
			// log.Println(err)
			return err
		}

		_, err = stmt.Exec(username, sid)
		if err != nil {
			// log.Println(err)
			return err
		}
	}

	return nil
}

func GetUserSessionData(db *sql.DB, store *session.Store, c *fiber.Ctx) (*models.User, error) {
	// Get current session
	s, _ := store.Get(c)
	// fmt.Println(s.Keys())

	// If there is a valid session
	if len(s.Keys()) > 0 {
		sid := s.ID()
		// From the session that is started we obtain the user id
		username := s.Get("username").(string)
		// Then with its uid we get the user data
		user := new(models.User)
		user.Username = username
		recoveredUser := st.GetUser(db, user.Username)

		// Get profile info
		U := &models.User{
			Email:    recoveredUser.Email,
			Username: recoveredUser.Username,
			Session:  sid,
		}

		// Get sessions list
		rows, err := db.Query(`SELECT v, e FROM session_store WHERE u = $1`, username)
		if err != nil {
			log.Println(err)
		}

		defer rows.Close()

		// Loop through sessions
		for rows.Next() {
			var (
				data       = []byte{}
				exp  int64 = 0
			)
			if err := rows.Scan(&data, &exp); err != nil {
				log.Println(err)
				return nil, err
			}

			// If session isn't expired
			if exp > time.Now().Unix() {
				// Decode session data
				gd := gob.NewDecoder(bytes.NewBuffer(data))
				dm := make(map[string]interface{})
				if err := gd.Decode(&dm); err != nil {
					log.Println(err)
					return nil, err
				}

				// Append session
				U.Sessions = append(
					U.Sessions,
					models.UserSession{
						SID:    dm["sid"].(string),
						IP:     dm["ip"].(string),
						Login:  dm["login"].(string),
						Expiry: time.Unix(exp, 0).UTC().String(),
						UA:     dm["ua"].(string),
					},
				)
			}
		}

		return U, nil
	}

	return nil, nil
}

func RemoveUserSession(store *session.Store, c *fiber.Ctx) (bool, error) {
	//Get session ID
	sid := c.Query("sid")
	// fmt.Println("SID: ", sid)

	// Get current session
	s, _ := store.Get(c)
	// fmt.Println(s.Fresh())

	// Check session ID
	if len(sid) > 0 {
		// Get requested session
		data, err := store.Storage.Get(sid)
		if err != nil {
			return false, err
		}

		// Decode requested session data
		gd := gob.NewDecoder(bytes.NewBuffer(data))
		dm := make(map[string]interface{})
		if err := gd.Decode(&dm); err != nil {
			return false, err
		}

		// If it belongs to current user destroy requested session
		if s.Get("username").(string) == dm["username"] {
			store.Storage.Delete(sid)
		}

		return false, nil
	} else {
		// Destroy current session
		s.Destroy()
	}

	return true, nil
}
