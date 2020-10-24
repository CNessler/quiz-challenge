import React from 'react';
import './App.css';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      items: []
    };
  }

    handleFormSubmit(e) {
        e.preventDefault();
        let userData = this.state.newUser;

        fetch("http://example.com", {
            method: "POST",
            body: JSON.stringify(userData),
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json"
            }
        }).then(response => {
            response.json().then(data => {
                console.log("Successful" + data);
            });
        });
    }

  componentDidMount() {
    fetch("/stars")
        .then(res => res.json())
        .then(
            (result) => {
              this.setState({
                isLoaded: true,
                items: result
              });
            },
            // Note: it's important to handle errors here
            // instead of a catch() block so that we don't swallow
            // exceptions from actual bugs in components.
            (error) => {
              this.setState({
                isLoaded: true,
                error
              });
            }
        )
  }

  render() {
    const { error, isLoaded, items } = this.state;
  //   if (error) {
  //     return <div>Error: {error.message}</div>;
  //   } else if (!isLoaded) {
  //     return <div>Loading...</div>;
  //   } else {
  //     return (
  //         <div>
  //           <h1>StarManager</h1>
  //           <ul>
  //             {items.map(item => (
  //                 <li key={item.id}>
  //                   <a href={item.url}>{item.name}</a> {item.description}
  //                 </li>
  //             ))}
  //           </ul>
  //         </div>
  //     );
  //   }
  // }
      return (

          <form className="container-fluid" onSubmit={this.handleFormSubmit}>
              <Input
                  inputType={"text"}
                  title={"Full Name"}
                  name={"name"}
                  value={this.state.newGuest.name}
                  placeholder={"Enter your name"}
                  handleChange={this.handleInput}
              />{" "}
          </form>
      )
}

export default App;