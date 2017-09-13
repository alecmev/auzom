import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link } from 'react-router';

import Button from '../components/Button';

export default class TeamUserTeam extends Component {
  static propTypes = {
    item: PropTypes.object.isRequired, // eslint-disable-line
    user: ImmutablePropTypes.map,
    isMe: PropTypes.bool,
    amLeader: PropTypes.bool,
    amAdmin: PropTypes.bool,
    leaderCount: PropTypes.number,
    patchUserTeam: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
    secondaryClass: PropTypes.string.isRequired,
  };

  handlePromote = () => {
    this.props.patchUserTeam(this.props.item.id, 'promote', () =>
      this.props.messagePush('promoted a member to leader'),
    );
  };

  handleDemote = () => {
    this.props.patchUserTeam(this.props.item.id, 'demote', () =>
      this.props.messagePush('demoted a leader'),
    );
  };

  handleKick = () => {
    const { user } = this.props;
    if (!window.confirm(`Kick ${user && user.get('nickname')}?`)) return;
    this.props.patchUserTeam(this.props.item.id, 'kick', () =>
      this.props.messagePush('kicked a member'),
    );
  };

  render() {
    const { item, user, isMe, amLeader, amAdmin, leaderCount } = this.props;
    const amCapable = amLeader || amAdmin;
    return (
      <li>
        {user && <Link to={`/users/${user.get('id')}`}>
          {user.get('nickname')}
        </Link>}
        <span className={this.props.secondaryClass}>
          {item.isLeader && ' / leader'}
        </span>
        {amCapable && !item.isLeader &&
          <Button
            text="promote"
            type="important"
            size="small"
            onClick={this.handlePromote}
          />
        }{amCapable && item.isLeader && leaderCount > 1 &&
          <Button
            text="demote"
            type="important"
            size="small"
            onClick={this.handleDemote}
          />
        }{!isMe && amCapable && (!item.isLeader || leaderCount > 1) &&
          <Button
            text="kick"
            type="important"
            size="small"
            onClick={this.handleKick}
          />
        }
      </li>
    );
  }
}
