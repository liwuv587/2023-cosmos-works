package keeper

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"planet/x/blog/types"
)

// TransmitUpdatePostPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitUpdatePostPacket(
	ctx sdk.Context,
	packetData types.UpdatePostPacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %w", err)
	}

	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvUpdatePostPacket processes packet reception
func (k Keeper) OnRecvUpdatePostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.UpdatePostPacketData) (packetAck types.UpdatePostPacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// TODO: packet reception logic
	id, err := strconv.ParseUint(data.GetPostID(), 10, 64)
	if err != nil {
		return packetAck, err
	}
	post, found := k.GetPost(ctx, id)
	if !found {
		return packetAck, errors.New("post not found")
	}
	post.Title = data.GetTitle()
	post.Content = data.GetContent()
	k.SetPost(ctx, post)
	packetAck.PostID = data.PostID

	return packetAck, nil
}

// OnAcknowledgementUpdatePostPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementUpdatePostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.UpdatePostPacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.UpdatePostPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic
		id, _ := strconv.ParseUint(data.PostID, 10, 64) // 字符串转换为 int
		post, found := k.GetSentPost(ctx, id)
		if !found {
			return errors.New("cannot find post")
		}
		post.Title = data.GetTitle()
		k.SetSentPost(ctx, post)
		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutUpdatePostPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutUpdatePostPacket(ctx sdk.Context, packet channeltypes.Packet, data types.UpdatePostPacketData) error {

	// TODO: packet timeout logic

	return nil
}
