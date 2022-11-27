// Code generated by smithy-go-codegen DO NOT EDIT.

package autoscaling

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Gets information about the scaling policies in the account and Region.
func (c *Client) DescribePolicies(ctx context.Context, params *DescribePoliciesInput, optFns ...func(*Options)) (*DescribePoliciesOutput, error) {
	if params == nil {
		params = &DescribePoliciesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DescribePolicies", params, optFns, c.addOperationDescribePoliciesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DescribePoliciesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DescribePoliciesInput struct {

	// The name of the Auto Scaling group.
	AutoScalingGroupName *string

	// The maximum number of items to be returned with each call. The default value is
	// 50 and the maximum value is 100.
	MaxRecords *int32

	// The token for the next set of items to return. (You received this token from a
	// previous call.)
	NextToken *string

	// The names of one or more policies. If you omit this property, all policies are
	// described. If a group name is provided, the results are limited to that group.
	// If you specify an unknown policy name, it is ignored with no error. Array
	// Members: Maximum number of 50 items.
	PolicyNames []string

	// One or more policy types. The valid values are SimpleScaling, StepScaling,
	// TargetTrackingScaling, and PredictiveScaling.
	PolicyTypes []string

	noSmithyDocumentSerde
}

type DescribePoliciesOutput struct {

	// A string that indicates that the response contains more items than can be
	// returned in a single response. To receive additional items, specify this string
	// for the NextToken value when requesting the next set of items. This value is
	// null when there are no more items to return.
	NextToken *string

	// The scaling policies.
	ScalingPolicies []types.ScalingPolicy

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDescribePoliciesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsquery_serializeOpDescribePolicies{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsquery_deserializeOpDescribePolicies{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDescribePolicies(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

// DescribePoliciesAPIClient is a client that implements the DescribePolicies
// operation.
type DescribePoliciesAPIClient interface {
	DescribePolicies(context.Context, *DescribePoliciesInput, ...func(*Options)) (*DescribePoliciesOutput, error)
}

var _ DescribePoliciesAPIClient = (*Client)(nil)

// DescribePoliciesPaginatorOptions is the paginator options for DescribePolicies
type DescribePoliciesPaginatorOptions struct {
	// The maximum number of items to be returned with each call. The default value is
	// 50 and the maximum value is 100.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribePoliciesPaginator is a paginator for DescribePolicies
type DescribePoliciesPaginator struct {
	options   DescribePoliciesPaginatorOptions
	client    DescribePoliciesAPIClient
	params    *DescribePoliciesInput
	nextToken *string
	firstPage bool
}

// NewDescribePoliciesPaginator returns a new DescribePoliciesPaginator
func NewDescribePoliciesPaginator(client DescribePoliciesAPIClient, params *DescribePoliciesInput, optFns ...func(*DescribePoliciesPaginatorOptions)) *DescribePoliciesPaginator {
	if params == nil {
		params = &DescribePoliciesInput{}
	}

	options := DescribePoliciesPaginatorOptions{}
	if params.MaxRecords != nil {
		options.Limit = *params.MaxRecords
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribePoliciesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribePoliciesPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribePolicies page.
func (p *DescribePoliciesPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*DescribePoliciesOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxRecords = limit

	result, err := p.client.DescribePolicies(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

func newServiceMetadataMiddleware_opDescribePolicies(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "autoscaling",
		OperationName: "DescribePolicies",
	}
}
