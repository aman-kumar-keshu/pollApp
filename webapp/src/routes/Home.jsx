import React, { useEffect, useState } from "react";
import PollList from "../components/PollList";
import Header from "../components/Header";
import { fetchPolls, updatePolls, deletePoll } from "../services/pollService";

const sortPolls = (polls) => polls.sort((poll1, poll2) => poll1.id - poll2.id);
function Home() {
  const [polls, setPolls] = useState([]);
  const [token, setToken] = useState("");

  useEffect(() => {
    const getData = async () => {
      const token = localStorage.getItem("token");
      const data = await fetchPolls();
      console.log("data list fetched from DB", data.polls);
      setPolls(sortPolls(data.polls.items));
      setToken(token);
    };
    getData();
  }, []);

  const handleUpdate = async (id, updatedPoll) => {
    const oldListWithoutId = polls.filter((poll) => poll.id !== id);
    const newPolls = sortPolls([...oldListWithoutId, updatedPoll]);

    setPolls(newPolls);
    updatePolls(id, updatedPoll);
  };

  const handleDelete = (id) => {
    const oldListWithoutId = polls.filter((poll) => poll.id !== id);
    setPolls(oldListWithoutId);
    deletePoll(id);
    console.log("handle delete");
  };

  const clearToken = () => {
    setToken("");
  };

  return (
    <>
      <Header numPolls={polls.length} token={token} clearToken={clearToken} />
      <div className="wrapper">
        <div className="card frame">
          <PollList
            polls={polls}
            onDelete={handleDelete}
            onUpdate={handleUpdate}
          />
        </div>
      </div>
    </>
  );
}

export default Home;
