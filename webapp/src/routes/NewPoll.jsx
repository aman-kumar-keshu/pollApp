import { useState } from "react";

function NewPoll() {
  const [enteredValue, setEnteredValue] = useState("");

  const handleNewPostSubmit = (e) => {
    e.preventDefault();
    console.log("add new poll callback", enteredValue);
  };

  return (
    <form onSubmit={handleNewPostSubmit}>
      <input
        type="text"
        value={enteredValue}
        onChange={(e) => setEnteredValue(e.target.value)}
      />
    </form>
  );
}

export default NewPoll;
