import React, { useEffect, useState } from "react";
import PollList from "../components/PollList";
import Header from "../components/Header";
import {
  fetchPolls,
  updatePolls,
  deletePoll,
  updateOption,
} from "../services/pollService";

const sort = (polls) => polls.sort((poll1, poll2) => poll1.id - poll2.id);
function Home() {
  const [polls, setPolls] = useState([]);
  const [token, setToken] = useState("");

  useEffect(() => {
    const getData = async () => {
      const token = localStorage.getItem("token");
      const data = await fetchPolls();
      console.log("data list fetched from DB", data.polls);
      for (var poll in data.polls.items) {
        sort(data.polls.items[poll].options);
      }

      setPolls(sort(data.polls.items));
      setToken(token);
    };
    getData();
  }, []);

  const handleUpdate = async (id, updatedOption) => {
    const pollId = updatedOption.PollId;
    const updatedPoll = polls.find((poll) => poll.id == pollId);
    const oldListWithoutId = polls.filter((poll) => poll.id !== pollId);
    // console.log("UpdatedPoll", updatedPoll, pollId); // 0

    const oldListWithoutOption = updatedPoll.options.filter(
      (option) => option.id !== updatedOption.id
    );
    const newOptions = sort([...oldListWithoutOption, updatedOption]);
    console.log("updated newOptions", newOptions);

    updatedPoll.options = newOptions;
    const newPolls = sort([...oldListWithoutId, updatedPoll]);
    setPolls(newPolls);
    updateOption(id, updatedOption);
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
