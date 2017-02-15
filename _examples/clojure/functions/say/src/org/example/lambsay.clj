(ns org.example.lambsay
    "What does the lamb say?"
    (:gen-class :implements [com.amazonaws.services.lambda.runtime.RequestStreamHandler])
    (:require [cheshire.core :as json]
              [clojure.java.io :as io]
              [clojure.string :as str])
    (:import (com.amazonaws.services.lambda.runtime Context)))

(defn it-says
    "The lamb says-"
    [what]
    (->
"                   _,._      ┌outline┐
               __.'   _)     │sayline│
              <_,)'.-\"a\\    /└outline┘
                /' (    \\  /
    _.-----..,-'   (`\"--^
   //              |  
  (|   `;      ,   |  
    \\   ;.----/  ,/ 
     ) // /   | |\\ \\
     \\ \\\\`\\   | |/ /
      \\ \\\\ \\  | |\\/"
        (str/replace "outline" (apply str (repeat (count what) "─")))
        (str/replace "sayline" what)))

(defn -handleRequest
    [this is os context]
    (let [w (io/writer os)]
    (-> (json/parse-stream (io/reader is))
        (get "say")
        (it-says)
        (str "\n")
        (->> (.write w)))
    (.flush w)))
