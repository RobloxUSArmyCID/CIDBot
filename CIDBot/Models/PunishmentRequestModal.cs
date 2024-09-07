using Discord.Interactions;
using Discord;

namespace CIDBot.Models
{
    public class PunishmentRequestModal : IModal
    {
        public string Title => "Punishment Request Submission";

        [InputLabel("Offender:")]
        [ModalTextInput("offender-username", TextInputStyle.Short, "Username", maxLength: 20)]
        public string? OffenderUsername { get; set; }

        [InputLabel("Ticket ID:")]
        [ModalTextInput("ticket-id", TextInputStyle.Short, "Ex: 4124")]
        public string? TicketId { get; set; }

        [InputLabel("Context:")]
        [ModalTextInput("punishment-context", TextInputStyle.Paragraph, "A quick brief about the case. Specify the punishment given.", minLength: 50)]
        public string? Context { get; set; }

        [InputLabel("Punishment Request URL:")]
        [ModalTextInput("punishment-request-url", TextInputStyle.Short, "https://docs.google.com/")]
        public string? PunishmentRequestUrl { get; set; }

        [InputLabel("Case File URL:")]
        [ModalTextInput("case-file-url", TextInputStyle.Short, "https://docs.google.com/")]
        public string? CaseFileUrl { get; set; }
    }
}
