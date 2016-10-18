import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link } from 'react-router';

import Button from '../components/Button';

export default class TournamentTeamSeasonRequest extends Component {
  static propTypes = {
    item: PropTypes.object.isRequired,
    team: ImmutablePropTypes.map,
    amAdmin: PropTypes.bool,
    canAccept: PropTypes.bool,
    patchTeamSeasonRequest: PropTypes.func.isRequired,
    loadTeamSeasons: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  handleAccept = () => {
    const { team } = this.props;
    if (!window.confirm(`Accept ${team && team.get('name')}?`)) return;
    this.props.patchTeamSeasonRequest(this.props.item.id, 'yes', (x) => {
      x.decision && this.props.loadTeamSeasons({
        teamId: x.teamId,
        seasonId: x.seasonId,
      });
      this.props.messagePush('accepted an application to a season');
    });
  };

  handleDecline = () => {
    const { team } = this.props;
    if (!window.confirm(`Decline ${team && team.get('name')}?`)) return;
    this.props.patchTeamSeasonRequest(this.props.item.id, 'no', () =>
      this.props.messagePush('declined an application to a season'),
    );
  };

  render() {
    const { team, amAdmin, canAccept } = this.props;
    return (
      <li>
        {team && <Link to={`/teams/${team.get('id')}`}>
          {team.get('name')}
        </Link>}
        {amAdmin && (
          <span>
            {canAccept &&
              <Button
                text="accept" type="important" size="small"
                onClick={this.handleAccept}
              />
            }
            <Button
              text="decline" type="important" size="small"
              onClick={this.handleDecline}
            />
          </span>
        )}
      </li>
    );
  }
}
