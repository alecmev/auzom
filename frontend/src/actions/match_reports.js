import callAPI from './api';
import { loadMatch } from './matches';
import { loadUser } from './users';

const matchReportLoaders = {
  matchId: loadMatch,
  agreedUponBy: loadUser,
  createdBy: loadUser,
};

export function createMatchReport(body, onSuccess) {
  return callAPI({
    url: '/match_reports',
    method: 'POST',
    body,
    type: 'MATCH_REPORT__CREATE',
    storage: 'matchReports',
    onSuccess,
  });
}

export function loadMatchReport(id, deps) {
  return callAPI({
    url: `/match_reports/${id}`,
    type: 'MATCH_REPORT__LOAD',
    meta: { id },
    storage: 'matchReports',
    loaders: matchReportLoaders,
    deps,
  });
}

matchReportLoaders._self = loadMatchReport;

export function loadMatchReports(filters, ids, deps) {
  return callAPI({
    url: '/match_reports',
    filters,
    type: 'MATCH_REPORTS__LOAD',
    storage: 'matchReports',
    loaders: matchReportLoaders,
    ids,
    deps,
  });
}

export function patchMatchReport(id, action, onSuccess) {
  return callAPI({
    url: `/match_reports/${id}`,
    method: 'PATCH',
    body: { action },
    type: 'MATCH_REPORT__PATCH',
    storage: 'matchReports',
    onSuccess,
  });
}
