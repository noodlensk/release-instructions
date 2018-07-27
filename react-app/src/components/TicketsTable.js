import React from 'react';
import axios from 'axios';
import {
    Form,
    FormGroup,
    Input,
    Label,
    Button,
    Row,
    Container,
    Col,
    Table,
    Modal,
    ModalHeader,
    ModalBody,
    ModalFooter
} from 'reactstrap';

class CommitItem extends React.Component {
    render() {
        return (
            <tr>
                <td> {this.props.message} </td>
                <td> {this.props.repo} </td>
            </tr>
        );
    }
}

class TicketItem extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            modal: false
        };

        this.toggle = this.toggle.bind(this);
    }

    toggle() {
        this.setState({
            modal: !this.state.modal
        });
    }

    render() {
        const commits = this.props.commits.map((commit, i) => {
            return (
                <CommitItem key={i} id={commit.Key} message={commit.Message} repo={commit.Repo.Name} />
            );
        });
        return (
            <tr>
                <td> {this.props.id} </td>
                <td> {this.props.name} </td>
                <td>
                    <a onClick={this.toggle}>{this.props.commits.length}</a>
                    <Modal size='lg' isOpen={this.state.modal} toggle={this.toggle} className={this.props.className}>
                        <ModalHeader toggle={this.toggle}>Commits</ModalHeader>
                        <ModalBody>
                            <Table striped>
                                <thead>
                                    <tr>
                                        <th>Message</th><th>RepoName</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {commits}
                                </tbody>
                            </Table>
                        </ModalBody>
                    </Modal>
                </td>
            </tr>
        );
    }
}

export class TicketsTable extends React.Component {
    constructor(props) {
        super(props);
    }
    render() {
        const tickets = this.props.tickets.map((ticket, i) => {
            return (
                <TicketItem key={i} id={ticket.Key} name={ticket.Name} commits={ticket.Commits != null ? ticket.Commits : []} />
            );
        });
        return (
            <div>
                <h3>Tickets</h3>
                <Table striped>
                    <thead>
                        <tr>
                            <th>Id</th><th>Name</th><th>Commits</th>
                        </tr>
                    </thead>
                    <tbody>
                        {tickets}
                    </tbody>
                </Table>
            </div>
        )
    }
}