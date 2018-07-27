import config from '../config';

import axios from 'axios';
import React from 'react';
import {
    Form,
    FormGroup,
    Label,
    Input,
    Button
} from 'reactstrap';

export class SearchForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            project: '',
            fixVersion: '',
        };

        this.handleInputChange = this.handleInputChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }
    handleInputChange(event) {
        const target = event.target;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        const name = target.name;

        this.setState({
            [name]: value
        });

    }

    handleSubmit(event) {
        this.serverRequest =
            axios
                .get(`${config.API_URL}/api/instructions?project=` + this.state.project + "&fixVersion=" + this.state.fixVersion)
                .then((result) => {
                    console.log(result)
                    this.props.updateGlobalState(result.data.Tickets, result.data.AffectedRepos)
                });
        event.preventDefault();
    }
    render() {
        return (
            <Form inline onSubmit={this.handleSubmit}>
                <FormGroup className="mb-2 mr-sm-2 mb-sm-0">
                    <Label for="searchProject" className="mr-sm-2">Project</Label>
                    <Input
                        type="text"
                        name="project"
                        id="searchProject"
                        placeholder="JIRA project"
                        value={this.state.project}
                        onChange={this.handleInputChange}
                    />
                </FormGroup>
                <FormGroup className="mb-2 mr-sm-2 mb-sm-0">
                    <Label for="searchFixVersion" className="mr-sm-2">FixVersion</Label>
                    <Input
                        type="text"
                        name="fixVersion"
                        id="searchFixVersion"
                        placeholder=""
                        value={this.state.fixVersion}
                        onChange={this.handleInputChange}
                    />
                </FormGroup>
                <Button>Search</Button>
            </Form>
        )
    }
}