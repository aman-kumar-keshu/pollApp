import { useState } from "react";
import { createPoll } from "../services/pollService";
import { useNavigate } from "react-router-dom";

function NewPoll(props) {
  const [name, setName] = useState("");
  const [src, setSrc] = useState("");
  const [topic, setTopic] = useState("");
  const navigator = useNavigate();

  const handleNewPostSubmit = async (e) => {
    e.preventDefault();
    console.log("add new poll callback", name, src, topic);
    const poll = { name, src, topic };
    await createPoll(poll);
    navigator("/");
  };

  return (
    <form onSubmit={handleNewPostSubmit}>
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <input
        type="text"
        value={topic}
        onChange={(e) => setTopic(e.target.value)}
      />
      <input type="text" value={src} onChange={(e) => setSrc(e.target.value)} />

      <button type="submit"> Submit </button>
    </form>
  );
}

export default NewPoll;
