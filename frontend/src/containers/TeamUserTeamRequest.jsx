import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link } from 'react-router';

import Button from '../components/Button';

export default class TeamUserTeamRequest extends Component {
  static propTypes = {
    item: PropTypes.object.isRequired, // eslint-disable-line
    user: ImmutablePropTypes.map,
    amLeader: PropTypes.bool,
    amAdmin: PropTypes.bool,
    patchUserTeamRequest: PropTypes.func.isRequired,
    loadUserTeams: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
    secondaryClass: PropTypes.string.isRequired,
  };

  handleAcceptApplication = () => {
    this.props.patchUserTeamRequest(this.props.item.id, 'yes', (x) => {
      x.decision && this.props.loadUserTeams({
        userId: x.userId,
        teamId: x.teamId,
      });
      this.props.messagePush('accepted an application to a team');
    });
  };

  handleDeclineApplication = () => {
    this.props.patchUserTeamRequest(this.props.item.id, 'no', () =>
      this.props.messagePush('declined an application to a team'),
    );
  };

  handleCancelInvite = () => {
    this.props.patchUserTeamRequest(this.props.item.id, 'no', () =>
      this.props.messagePush('cancelled an invite to a team'),
    );
  };

  handleApprove = () => {
    this.props.patchUserTeamRequest(this.props.item.id, 'yes', (x) => {
      x.decision && this.props.loadUserTeams({
        userId: x.userId,
        teamId: x.teamId,
      });
      this.props.messagePush('approved a user-team request');
    });
  };

  handleDeny = () => {
    this.props.patchUserTeamRequest(this.props.item.id, 'no', () =>
      this.props.messagePush('denied a user-team request'),
    );
  };

  render() {
    const { item, user, amLeader, amAdmin } = this.props;
    return (
      <li>
        {user && <Link to={`/users/${user.get('id')}`}>
          {user.get('nickname')}
        </Link>}
        <span className={this.props.secondaryClass}>
          {item.userDecision && !item.userDecidedAt &&
            ' / application'
          }
          {item.leaderDecision && !item.leaderDecidedAt &&
            ' / invite'
          }
          {item.userDecision && item.leaderDecision &&
            ' / pending admins\' approval'
          }
        </span>
        {amLeader && item.userDecision && !item.leaderDecision &&
          <Button
            text="accept"
            type="important"
            size="small"
            onClick={this.handleAcceptApplication}
          />
        }{amLeader && item.userDecision &&
          <Button
            text="decline"
            type="important"
            size="small"
            onClick={this.handleDeclineApplication}
          />
        }{amLeader && item.leaderDecision && !item.userDecision &&
          <Button
            text="cancel invite"
            type="important"
            size="small"
            onClick={this.handleCancelInvite}
          />
        }{amAdmin && item.adminDecision === null &&
          <Button
            text="approve"
            type="important"
            size="small"
            onClick={this.handleApprove}
          />
        }{amAdmin && (item.adminDecision === null || item.adminDecidedAt) &&
          <Button
            text="deny"
            type="important"
            size="small"
            onClick={this.handleDeny}
          />
        }
      </li>
    );
  }
}
