import React from 'react';
import {
    Row,
    Col,
} from 'reactstrap';

import { SearchForm } from './SearchForm.js';
import { ReleaseInstructions } from './ReleaseInstructions.js';
import { TicketsTable } from './TicketsTable.js';



export class SearchPage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tickets: [],
            repos: []
        };
        this._updateGlobalState = this._updateGlobalState.bind(this)
    }

    _updateGlobalState(tickets, repos) {
        this.setState({
            tickets,
            repos
        })
    }

    render() {

        return (
            <div>
                <SearchForm updateGlobalState={this._updateGlobalState} />
                <br />
                {this.state.repos.length > 0 &&
                    <Row>
                        <Col>
                            <ReleaseInstructions repos={this.state.repos} />
                        </Col>
                        <Col>
                            <TicketsTable tickets={this.state.tickets} />
                        </Col>
                    </Row>}
            </div>
        );
    }
}