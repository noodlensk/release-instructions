import React from 'react';
import {
    Container,
} from 'reactstrap';
class ReleaseInstructionsItem extends React.Component {
    render() {
        return (
            <li>{this.props.step}.) Update {this.props.repo.Name}</li>
        )
    }
}
export class ReleaseInstructions extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const instructions = this.props.repos.sort((a, b) => a.Weight > b.Weight).map((repo, i) => {
            return (
                <ReleaseInstructionsItem key={i} step={i + 1} repo={repo} />
            );
        });
        return (
            <Container>
                <h4>Instructions</h4>
                <ul>
                    {instructions}
                </ul>
            </Container>
        )
    }
}