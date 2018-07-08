const backendapi = require('../data/backendapi');

const maxPercentage = 100;

class BuildListener
{
    /**
     * @param {context.Context} context
     * @param {Request} request
     * @param {Response} response
     */
    constructor(context)
    {
        /**
         * @property {context.Context}
         */
        this._context = context;
        this._backendApi = new backendapi.BackendApi(this._context.config);
    }

    onBuildFinished(uuid, succeed)
    {
        const status = succeed ? 'succeed' : 'failed';
        console.log(`BuildListener: build ${uuid} ${status}`);

        const promise = this._processBuild(uuid, succeed);
        promise.catch((error) => {
            console.error('failed to process finished build', error);
        });
    }

    async _processBuild(uuid, succeed)
    {
        const repository = this._context.connectDB();
        let score = 0;
        if (succeed)
        {
            const report = await this._backendApi.getBuildReport(uuid);
            if (report.tests_total > 0)
            {
                score = maxPercentage * report.tests_passed / report.tests_total;
            }
        }
        const commitInfo = await repository.getCommitInfo(uuid);
        const solutionId = parseInt(commitInfo['solution_id']);
        await repository.updateCommit(uuid, succeed, score);
        await repository.updateSolution(solutionId, score);
    }
}

module.exports.BuildListener = BuildListener;
