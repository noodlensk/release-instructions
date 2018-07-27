// -*- JavaScript -*-
import 'bootstrap/dist/css/bootstrap.min.css';
import React from 'react';
import ReactDOM from 'react-dom';
import {
  Navbar,
  NavbarBrand,
  Container
} from 'reactstrap';
import { SearchPage } from './components/SearchPage';

class App extends React.Component {
  render() {
    return (
      <div>
        <Navbar color="light" light >
          <NavbarBrand href="/">Release instructions</NavbarBrand>
        </Navbar>
        <br />
        <Container>
          <SearchPage />
        </Container>
      </div>
    );
  }
}

ReactDOM.render(<App />, document.querySelector("#app"));
module.hot.accept();
