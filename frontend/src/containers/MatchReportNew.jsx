import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

import Button from '../components/Button';
import Checkbox from '../components/Checkbox';
import Input from '../components/Input';

const swapSidesEveryRounds = 1;

@connect(createStructuredSelector({
  gameMaps: selectors.gameMaps,
  teams: selectors.teams,
  amAdmin: selectors.amAdmin,
}), actions)
export default class MatchReportNew extends Component {
  static propTypes = {
    bracket: ImmutablePropTypes.map,
    matchId: PropTypes.string,
    matchPath: PropTypes.string,
    match: ImmutablePropTypes.map,
    matchMaps: ImmutablePropTypes.list,
    amAdmin: PropTypes.bool,

    gameMaps: ImmutablePropTypes.map,
    teams: ImmutablePropTypes.map,

    createMatchReport: PropTypes.func.isRequired,
    loadMatch: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  state = {};

  submit = () => {
    const {
      bracket, matchId, matchPath, matchMaps, messagePush, createMatchReport,
    } = this.props;
    if (!bracket || !matchId || !matchMaps || !matchMaps.size) return;
    const body = {
      matchId,
      rounds: [],
    };
    const roundsPerMap = (bracket.get('type') === 'ace-pre-swiss') ? 1 : 2;
    for (const map of matchMaps.filter(x => !x.get('isBan')).toJS()) {
      for (
        let i = 0, swap = false;
        i < roundsPerMap;
        i += 1, swap = (i % swapSidesEveryRounds) ? swap : !swap
      ) {
        const key = `${map.id}-${i}`;
        const round = {
          gameMapId: map.gameMapId,
          isTeamXOnSideY: swap,
          isNotPlayed: this.r[`${key}-notPlayed`].state.value,
          rawScoreX: this
            .r[`${key}-rawScoreX`]
            .state
            .value
            .trim() || '0',
          rawScoreY: this
            .r[`${key}-rawScoreY`]
            .state
            .value
            .trim() || '0',
        };
        if (
          !round.isNotPlayed &&
          (round.rawScoreX === '' || round.rawScoreY === '')
        ) {
          messagePush('some scores missing', true);
          return;
        }

        if (this.state[key]) {
          round.penalties = [];
          for (let j = 0; j < this.state[key]; j += 1) {
            round.penalties.push({
              reason: this.r[`${key}-penaltyReason-${j}`].state.value.trim(),
              rawScoreX: this.r[`${key}-rawScoreXPenalty-${j}`].state.value
                .trim() || '0',
              rawScoreY: this.r[`${key}-rawScoreYPenalty-${j}`].state.value
                .trim() || '0',
            });
          }
        }

        if (this.state[`${key}-override`]) {
          round.overrideReason = this.r[`${key}-overrideReason`].state.value
            .trim();
          round.isPenalOverride = this.r[`${key}-isPenalOverride`].state
            .value;
          round.rawScoreXOverride = this.r[`${key}-rawScoreXOverride`].state
            .value.trim() || null;
          round.rawScoreYOverride = this.r[`${key}-rawScoreYOverride`].state
            .value.trim() || null;
        }

        body.rounds.push(round);
      }
    }

    if (this.state.wholeMatch) {
      body.penalties = [];
      for (let i = 0; i < this.state.wholeMatch; i += 1) {
        body.penalties.push({
          reason: this.r[`wholeMatch-penaltyReason-${i}`].state.value.trim(),
          scoreX: this.r[`wholeMatch-scoreXPenalty-${i}`].state.value
            .trim() || '0',
          scoreY: this.r[`wholeMatch-scoreYPenalty-${i}`].state.value
            .trim() || '0',
          rawScoreX: this.r[`wholeMatch-rawScoreXPenalty-${i}`].state.value
            .trim() || '0',
          rawScoreY: this.r[`wholeMatch-rawScoreYPenalty-${i}`].state.value
            .trim() || '0',
        });
      }
    }

    if (this.state['wholeMatch-override']) {
      body.overrideReason = this.r['wholeMatch-overrideReason'].state.value
        .trim();
      body.isPenalOverride = this.r['wholeMatch-isPenalOverride'].state
        .value;
      body.scoreXOverride = this.r['wholeMatch-scoreXOverride'].state.value
        .trim() || null;
      body.scoreYOverride = this.r['wholeMatch-scoreYOverride'].state.value
        .trim() || null;
      body.rawScoreXOverride = this.r['wholeMatch-rawScoreXOverride'].state
        .value.trim() || null;
      body.rawScoreYOverride = this.r['wholeMatch-rawScoreYOverride'].state
        .value.trim() || null;
    }

    createMatchReport(body, () => {
      if (this.props.amAdmin) {
        this.props.loadMatch(matchId);
        matchPath && this.context.router.push(`${matchPath}/details`);
      } else {
        matchPath && this.context.router.push(`${matchPath}/reports`);
      }

      messagePush('submitted a match report');
    });
  };

  handleToggleOverride = (_, key) => {
    this.setState({ [key]: !this.state[key] });
  };

  handleAddPenalty = (_, key) => {
    this.setState({ [key]: this.state[key] + 1 || 1 });
  };

  handleRemovePenalty = (_, key) => {
    this.setState({ [key]: this.state[key] ? this.state[key] - 1 : 0 });
  };

  render() {
    const { bracket, match, matchMaps, teams, gameMaps, amAdmin } = this.props;
    if (
      !bracket || !match || !match.get('areMapsReady') || !matchMaps ||
      !matchMaps.size
    ) return null;
    const roundsPerMap = (bracket.get('type') === 'ace-pre-swiss') ? 1 : 2;
    return (
      <div className={this.cni()}>
        <div className={this.cn({ d: 'inner' })}>
          {matchMaps.filter(x => !x.get('isBan')).toJS().map((x) => {
            const map = gameMaps.get(x.gameMapId);
            let mapName;
            let sideX;
            let sideY;
            if (map) {
              [mapName, sideX, sideY] = [
                map.get('abbr'), map.get('sideXAbbr'), map.get('sideYAbbr'),
              ];
            }

            const res = [];
            for (
              let i = 0, swap = false;
              i < roundsPerMap;
              i += 1, swap = (i % swapSidesEveryRounds) ? swap : !swap
            ) {
              const key = `${x.id}-${i}`;
              const [teamX, teamY] = [
                teams.get(match.get('teamX')),
                teams.get(match.get('teamY')),
              ];
              res.push(
                <div
                  key={key}
                  className={this.cn({ d: 'round' })}
                >
                  <div className={this.cn({ d: 'label' })}>
                    <div className={this.cn({ d: 'labelMap' })}>{mapName}</div>
                    <div className={this.cn({ d: 'labelRound' })}>
                      Round {i + 1}
                    </div>
                  </div>
                  <Input
                    ref={this.rcb(`${key}-rawScoreX`)}
                    className={this.cn({ d: 'score', m: 'left' })}
                    label={`${
                      swap ? sideY : sideX}${
                      teamX ? ` / ${teamX.get('abbr')}` : ''
                    }`}
                  />
                  <Input
                    ref={this.rcb(`${key}-rawScoreY`)}
                    className={this.cn({ d: 'score' })}
                    label={`${
                      swap ? sideX : sideY}${
                      teamY ? ` / ${teamY.get('abbr')}` : ''
                    }`}
                  />
                  <Checkbox
                    ref={this.rcb(`${key}-notPlayed`)}
                    className={this.cn({ d: 'checkbox' })}
                    label="Wasn't played"
                  />
                  {amAdmin && this.state[`${key}-override`] && <div>
                    <Input
                      ref={this.rcb(`${key}-overrideReason`)}
                      className={this.cn({ d: 'reason' })}
                      label="Override reason"
                    />
                    <Checkbox
                      ref={this.rcb(`${key}-isPenalOverride`)}
                      className={this.cn({ d: 'checkbox' })}
                      label="Is penal override"
                    />
                    <Input
                      ref={this.rcb(`${key}-rawScoreXOverride`)}
                      className={this.cn({ d: 'score', m: 'left' })}
                      label={`${swap ? sideY : sideX} override`}
                    />
                    <Input
                      ref={this.rcb(`${key}-rawScoreYOverride`)}
                      className={this.cn({ d: 'score' })}
                      label={`${swap ? sideX : sideY} override`}
                    />
                  </div>}
                  {amAdmin && <Button
                    block
                    className={this.cn({ d: 'button' })}
                    type="important"
                    text="toggle override"
                    size="small"
                    meta={`${key}-override`}
                    onClick={this.handleToggleOverride}
                  />}
                  {amAdmin && !!this.state[key] &&
                  [...Array(this.state[key])].map((_, j) =>
                    <div key={`${key}-penalty-${j}`}>
                      <Input
                        ref={this.rcb(`${key}-penaltyReason-${j}`)}
                        className={this.cn({ d: 'reason' })}
                        label={`Penalty ${j} reason`}
                      />
                      <Input
                        ref={this.rcb(`${key}-rawScoreXPenalty-${j}`)}
                        className={this.cn({ d: 'score', m: 'left' })}
                        label={`${swap ? sideY : sideX} penalty ${j}`}
                      />
                      <Input
                        ref={this.rcb(`${key}-rawScoreYPenalty-${j}`)}
                        className={this.cn({ d: 'score' })}
                        label={`${swap ? sideX : sideY} penalty ${j}`}
                      />
                    </div>
                  )}
                  {amAdmin && <Button
                    block
                    className={this.cn({ d: 'button' })}
                    type="important"
                    text="add penalty"
                    size="small"
                    meta={key}
                    onClick={this.handleAddPenalty}
                  />}
                  {amAdmin && !!this.state[key] && <Button
                    block
                    className={this.cn({ d: 'button' })}
                    type="important"
                    text="remove penalty"
                    size="small"
                    meta={key}
                    onClick={this.handleRemovePenalty}
                  />}
                </div>
              );
            }

            return res;
          })}
          {amAdmin && <div className={this.cn({ d: 'round' })}>
            <div className={this.cn({ d: 'label' })}>
              <div className={this.cn({ d: 'labelMap' })}>whole match</div>
            </div>
            {this.state['wholeMatch-override'] && <div>
              <Input
                ref={this.rcb('wholeMatch-overrideReason')}
                className={this.cn({ d: 'reason' })}
                label="Override reason"
              />
              <Checkbox
                ref={this.rcb('wholeMatch-isPenalOverride')}
                className={this.cn({ d: 'checkbox' })}
                label="Is penal override"
              />
              <Input
                ref={this.rcb('wholeMatch-scoreXOverride')}
                className={this.cn({ d: 'score', m: 'left' })}
                label="X override"
              />
              <Input
                ref={this.rcb('wholeMatch-scoreYOverride')}
                className={this.cn({ d: 'score' })}
                label="Y override"
              />
              <div />
              <Input
                ref={this.rcb('wholeMatch-rawScoreXOverride')}
                className={this.cn({ d: 'score', m: 'left' })}
                label="Raw X override"
              />
              <Input
                ref={this.rcb('wholeMatch-rawScoreYOverride')}
                className={this.cn({ d: 'score' })}
                label="Raw Y override"
              />
            </div>}
            <Button
              block
              className={this.cn({ d: 'button' })}
              type="important"
              text="toggle override"
              size="small"
              meta="wholeMatch-override"
              onClick={this.handleToggleOverride}
            />
            {!!this.state.wholeMatch &&
            [...Array(this.state.wholeMatch)].map((_, i) =>
              <div key={`wholeMatch-penalty-${i}`}>
                <Input
                  ref={this.rcb(`wholeMatch-penaltyReason-${i}`)}
                  className={this.cn({ d: 'reason' })}
                  label={`Penalty ${i} reason`}
                />
                <Input
                  ref={this.rcb(`wholeMatch-scoreXPenalty-${i}`)}
                  className={this.cn({ d: 'score', m: 'left' })}
                  label={`X penalty ${i}`}
                />
                <Input
                  ref={this.rcb(`wholeMatch-scoreYPenalty-${i}`)}
                  className={this.cn({ d: 'score' })}
                  label={`Y penalty ${i}`}
                />
                <div />
                <Input
                  ref={this.rcb(`wholeMatch-rawScoreXPenalty-${i}`)}
                  className={this.cn({ d: 'score', m: 'left' })}
                  label={`Raw X penalty ${i}`}
                />
                <Input
                  ref={this.rcb(`wholeMatch-rawScoreYPenalty-${i}`)}
                  className={this.cn({ d: 'score' })}
                  label={`Raw Y penalty ${i}`}
                />
              </div>
            )}
            <Button
              block
              className={this.cn({ d: 'button' })}
              type="important"
              text="add penalty"
              size="small"
              meta="wholeMatch"
              onClick={this.handleAddPenalty}
            />
            {!!this.state.wholeMatch && <Button
              block
              className={this.cn({ d: 'button' })}
              type="important"
              text="remove penalty"
              size="small"
              meta="wholeMatch"
              onClick={this.handleRemovePenalty}
            />}
          </div>}
          <Button
            block
            className={this.cn({ d: 'button' })}
            type="important"
            text="submit report"
            onClick={this.submit}
          />
        </div>
      </div>
    );
  }
}
