import "./TodoList.scss";

export interface TodoListEntry {
  id: number;
  description: string;
  done: boolean;
}

interface TodoListProps {
  listEntries: TodoListEntry[];
}

export function TodoList(props: TodoListProps) {
    console.log(props)
  return (
    <>
      {props.listEntries.map((entry) => {
        console.log(entry)
        return (
          <div key={entry.id} className="todo-list-entry">
            <p>{entry.description}</p>
            <input type="checkbox"></input>
          </div>
        );
      })}
    </>
  );
}
