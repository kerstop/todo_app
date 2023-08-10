import "./TodoList.scss";

export interface TodoListEntry {
  id: number;
  description: string;
  complete: boolean;
}

interface TodoListProps {
  listEntries: TodoListEntry[];
  setListEntries: (newList: TodoListEntry[]) => void;
}

export function TodoList(props: TodoListProps) {
  const toggleDone = (id: number) => {
    fetch("/api/todo/toggleDone", {
      method: "POST",
      body: JSON.stringify(id),
    });

    const newList = [...props.listEntries];
    const this_element = newList.find((v) => v.id === id);
    if (this_element !== undefined) {
      this_element.complete = !this_element.complete;
    }
    props.setListEntries(newList)
  };

  return (
    <>
      {props.listEntries.map((entry) => {
        console.log(entry);
        return (
          <div key={entry.id} className="todo-list-entry">
            <p>{entry.description}</p>
            <input
              type="checkbox"
              checked={entry.complete}
              onClick={() => {
                toggleDone(entry.id);
              }}
            ></input>
          </div>
        );
      })}
    </>
  );
}
