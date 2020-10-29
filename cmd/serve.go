/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run http service",
	Run: func(cmd *cobra.Command, args []string) {
		listen := ":3000"
		router := mux.NewRouter()

		router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, I'm buhaoyong\n"))
		})

		router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
		router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
		router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("POST")
		router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("PUT")
		router.HandleFunc("/user/{id}", func(w http.ResponseWriter, r *http.Request) {}).Methods("DELETE")

		server := &http.Server{
			Handler: router,
			Addr:    listen,
		}

		go func() {
			log.Println("Start server on " + listen)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		log.Println("Shutdown...")

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
