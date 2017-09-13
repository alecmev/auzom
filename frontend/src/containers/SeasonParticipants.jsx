import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link } from 'react-router';

import Component from '../utils/Component';

import Markdown from '../components/Markdown';
import Img from '../components/Img';

export default class SeasonParticipants extends Component {
  static propTypes = {
    location: PropTypes.object, // eslint-disable-line
    participants: ImmutablePropTypes.list,
    applications: ImmutablePropTypes.list,
    teams: ImmutablePropTypes.map,
  };

  render() {
    const { participants, applications, teams } = this.props;
    if (!participants || !applications) return null;

    const participantsItems = [];
    participants.toJS().forEach((x) => {
      const team = teams.get(x.teamId);
      participantsItems.push(
        <Link
          key={x.id}
          className={this.cn({
            d: 'item',
            m: x.leftAt && !x.isDone && 'fade',
          })}
          to={`${this.props.location.pathname}/${x.teamId}`}
        >
          {team && <div className={this.cn({ d: 'logoOuter' })}>
            {team.get('logo') && <Img
              className={this.cn({ d: 'logo' })}
              src={team.get('logo')}
            />}
          </div>}
          {team && <div className={this.cn({ d: 'name' })}>
            <div className={this.cn({ d: 'nameTag' })}>
              {team.get('abbr')}
            </div>
            <div className={this.cn({ d: 'nameFull' })}>
              {team.get('name')}
            </div>
          </div>}
          {x.leftAt && <div className={this.cn({ d: 'done' })}>
            {x.isDone ? '✔' : '❗'}
          </div>}
        </Link>,
      );
    });

    const applicationsItems = [];
    applications.filter(y => y.get('decision') === null).toJS().forEach((x) => {
      const team = teams.get(x.teamId);
      applicationsItems.push(
        <Link
          key={x.id}
          className={this.cn({ d: 'item' })}
          to={`${this.props.location.pathname}/${x.teamId}`}
        >
          {team && <div className={this.cn({ d: 'logoOuter' })}>
            {team.get('logo') && <Img
              className={this.cn({ d: 'logo' })}
              src={team.get('logo')}
            />}
          </div>}
          {team && <div className={this.cn({ d: 'name' })}>
            <div className={this.cn({ d: 'nameTag' })}>
              {team.get('abbr')}
            </div>
            <div className={this.cn({ d: 'nameFull' })}>
              {team.get('name')}
            </div>
          </div>}
        </Link>,
      );
    });

    // `https://unsplash.it/48/48/?image=${x.teamId}`

    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          <div className={this.cn({ d: 'inner' })}>
            {!!participantsItems.length && <div>
              <div className={this.cn({ d: 'heading' })}>
                <Markdown
                  className={this.cn({ d: 'headingTitle' })}
                  source="# Participants"
                />
              </div>
              <div className={this.cn({ d: 'items' })}>
                {participantsItems}
              </div>
            </div>}
            {!!applicationsItems.length && <div>
              <div className={this.cn({ d: 'heading' })}>
                <Markdown
                  className={this.cn({ d: 'headingTitle' })}
                  source="# Applicants"
                />
              </div>
              <div className={this.cn({ d: 'items' })}>
                {applicationsItems}
              </div>
            </div>}
          </div>
        </div>
      </div>
    );
  }
}
