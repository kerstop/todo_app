import * as React from "react";
import "./App.css";
import { TodoList, TodoListEntry } from "./TodoList";
import { LoginBar } from "./LoginBar";

function App() {
  const [list, setList] = React.useState<TodoListEntry[]>([]);
  const description_field = React.useRef<null| HTMLInputElement>(null)

  React.useEffect(() => {
    fetch(`/api/todo`)
      .then((val) => {
        console.log(val.body);
        return val.json();
      })
      .then((val) => {
        setList(val);
      });
  }, []);

  const createTodo: React.FormEventHandler = (e) => {
    e.preventDefault()

    if (description_field.current !== null) {
      const description = description_field.current.value
      if (description === null) return;
      console.log(description)
      fetch("/api/todo", {
        body: JSON.stringify({
          description: description,
        }),
        method: "POST"
      })
      .then((resp) => {
        return resp.json()
      }).then((resp: number) => {

        setList([...list, {
          id: resp,
          description: description,
          done: false
        }])
      })
    }
  };

  return (
    <>
      <LoginBar></LoginBar>
      <TodoList listEntries={list}></TodoList>
      <form onSubmit={createTodo}>
        <input ref={description_field} type="text" name="description" placeholder="Todo" required></input>
        <input type="button" value="Submit"></input>
      </form>
    </>
  );
}

export default App;
