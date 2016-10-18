import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link } from 'react-router';

import Button from '../components/Button';

export default class TournamentTeamSeason extends Component {
  static propTypes = {
    item: PropTypes.object.isRequired,
    team: ImmutablePropTypes.map,
    amAdmin: PropTypes.bool,
    patchTeamSeason: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
    doneClass: PropTypes.string.isRequired,
  };

  handleDone = () => {
    const { team } = this.props;
    if (!window.confirm(`Mark ${team && team.get('name')} as done?`)) return;
    this.props.patchTeamSeason(this.props.item.id, 'done', () =>
      this.props.messagePush('marked a team as done'),
    );
  };

  handleLeave = () => {
    const { team } = this.props;
    if (!window.confirm(`Mark ${team && team.get('name')} as left?`)) return;
    this.props.patchTeamSeason(this.props.item.id, 'leave', () =>
      this.props.messagePush('marked a team as left'),
    );
  };

  handleKick = () => {
    const { team } = this.props;
    if (!window.confirm(`Kick ${team && team.get('name')}?`)) return;
    this.props.patchTeamSeason(this.props.item.id, 'kick', () =>
      this.props.messagePush('kicked a team from a season'),
    );
  };

  render() {
    const { item, team, amAdmin } = this.props;
    return (
      <li>
        {team && <Link to={`/teams/${team.get('id')}`}>
          {team.get('name')}
        </Link>}
        {item.isDone &&
          <span className={this.props.doneClass}>
            done
          </span>
        }
        {amAdmin && !item.leftAt && (
          <span>
            <Button
              text="done" type="important" size="small"
              onClick={this.handleDone}
            />
            <Button
              text="leave" type="important" size="small"
              onClick={this.handleLeave}
            />
            <Button
              text="kick" type="important" size="small"
              onClick={this.handleKick}
            />
          </span>
        )}
      </li>
    );
  }
}
